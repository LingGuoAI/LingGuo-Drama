package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"spiritFruit/app/services"
	"spiritFruit/pkg/upload"
	"strings"

	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/characters" // 或者是角色表
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
)

func HandleGenerateImage(ctx context.Context, t *asynq.Task) error {
	var p myAsynq.GenerateImagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 1. 获取任务
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		console.Error(fmt.Sprintf("Task %d not found in DB", p.AsyncTaskID))
		return nil
	}

	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始生成角色图片", p.AsyncTaskID))

	// 2. 初始化 AI
	taskModel.UpdateProgress(20)

	// ==========================================
	// 🔴 生图任务优先查库，后查 .env
	// ==========================================

	// 1) 先用 .env 环境变量进行基础兜底初始化
	aiConfig := openai.Config{
		Provider: config.GetString("ai.provider", "openai"),

		// OpenAI 配置
		OpenAIBaseURL:    config.GetString("ai.openai.base_url"),
		OpenAIKey:        config.GetString("ai.openai.api_key"),
		OpenAIModel:      config.GetString("ai.openai.model"),
		OpenAIImageModel: config.GetString("ai.openai.image_model", "dall-e-3"),

		// GetGoAPI 配置 (新增)
		GetGoAPIBaseURL:    config.GetString("ai.getgoapi.base_url"),
		GetGoAPIKey:        config.GetString("ai.getgoapi.api_key"),
		GetGoAPIModel:      config.GetString("ai.getgoapi.model"),
		GetGoAPIImageModel: config.GetString("ai.getgoapi.image_model", "gpt-4o-image"),

		// Gemini 配置
		GeminiBaseURL: config.GetString("ai.gemini.base_url"),
		GeminiKey:     config.GetString("ai.gemini.api_key"),
		GeminiModel:   config.GetString("ai.gemini.model"),

		// 豆包 (Volcengine) 配置
		DoubaoBaseURL:    config.GetString("ai.doubao.base_url"),
		DoubaoKey:        config.GetString("ai.doubao.api_key"),
		DoubaoModel:      config.GetString("ai.doubao.model"),
		DoubaoImageModel: config.GetString("ai.doubao.image_model"),

		// Vertex AI 配置
		VertexKey:        config.GetString("ai.vertex.api_key"),
		VertexModel:      config.GetString("ai.vertex.model"),
		VertexImageModel: config.GetString("ai.vertex.image_model"),
	}

	// 2) 尝试从数据库加载优先级最高的 image (图片) 配置
	aiService := new(services.AiConfigService)
	errConfig, dbConfig := aiService.GetActiveConfigByType("image", taskModel.AdminID)

	if errConfig == nil && dbConfig.ID > 0 {
		providerName := *dbConfig.Provider
		baseURL := *dbConfig.BaseUrl
		apiKey := *dbConfig.ApiKey

		// 取 JSON 数组配置的第一个模型作为生图模型
		modelName := ""
		if len(dbConfig.Model) > 0 {
			modelName = dbConfig.Model[0]
		}

		switch providerName {
		case "getgoapi":
			aiConfig.Provider = "getgoapi"
			aiConfig.GetGoAPIBaseURL = baseURL
			aiConfig.GetGoAPIKey = apiKey
			if modelName != "" {
				aiConfig.GetGoAPIImageModel = modelName
			}

		case "openai":
			aiConfig.Provider = "openai"
			aiConfig.OpenAIBaseURL = baseURL
			aiConfig.OpenAIKey = apiKey
			if modelName != "" {
				aiConfig.OpenAIImageModel = modelName
			}

		case "gemini", "google":
			aiConfig.Provider = "gemini"
			aiConfig.GeminiBaseURL = baseURL
			aiConfig.GeminiKey = apiKey
			if modelName != "" {
				aiConfig.GeminiModel = modelName
			}

		case "doubao", "volcengine", "volces":
			aiConfig.Provider = "doubao"
			aiConfig.DoubaoBaseURL = baseURL
			aiConfig.DoubaoKey = apiKey
			if modelName != "" {
				aiConfig.DoubaoImageModel = modelName // 赋值给豆包图片模型字段
			}

		case "vertex":
			aiConfig.Provider = "vertex"
			aiConfig.VertexKey = apiKey
			if modelName != "" {
				aiConfig.VertexImageModel = modelName // 赋值给 Vertex 图片模型字段
			}

		default:
			// 兜底逻辑
			aiConfig.Provider = "openai"
			aiConfig.OpenAIBaseURL = baseURL
			aiConfig.OpenAIKey = apiKey
			if modelName != "" {
				aiConfig.OpenAIImageModel = modelName
			}
		}

		console.Success(fmt.Sprintf("任务[%d] - 成功挂载数据库 AI 图片配置: Provider=%s, Model=%s", p.AsyncTaskID, providerName, modelName))
	} else {
		console.Warning(fmt.Sprintf("任务[%d] - 未命中数据库 AI 图片配置，将降级使用 .env 默认配置", p.AsyncTaskID))
	}
	aiProvider := openai.NewProvider(aiConfig)

	// 3. 请求 AI 生图
	taskModel.UpdateProgress(40)
	req := openai.ImageRequest{
		Prompt: p.Prompt,
		N:      1,
		Size:   "1024x1024",
	}

	urls, err := aiProvider.GenerateImage(req)
	if err != nil {
		console.Error(fmt.Sprintf("AI 生图失败: %v", err))
		taskModel.MarkAsFailed(err)
		return err
	}

	if len(urls) == 0 {
		taskModel.MarkAsFailed(fmt.Errorf("no images generated"))
		return nil
	}

	rawImageURL := urls[0] // 可能是 http 链接，也可能是 data:image/png;base64...

	// 4. 保存到本地 (核心修改逻辑)
	taskModel.UpdateProgress(70)

	var localPath string
	var saveErr error

	if strings.HasPrefix(rawImageURL, "data:image") {
		//如果是 Base64 (Gemini)
		localPath, saveErr = upload.SaveBase64Image(rawImageURL)
	} else {
		// 如果是 URL (OpenAI)
		localPath, saveErr = upload.DownloadAndSave(rawImageURL)
	}

	if saveErr != nil {
		taskModel.MarkAsFailed(fmt.Errorf("save image locally failed: %v", saveErr))
		return saveErr
	}

	// 拼接完整的访问 URL (假设你的静态资源挂载在 /uploads 下)
	// 注意：这里取决于你的 Nginx 或 Gin Static 配置
	// 如果 localPath 是 "uploads/images/xxx.png"，且 Gin 路由是 r.Static("/uploads", "./uploads")
	// 那么前端访问路径就是 "/uploads/images/xxx.png"
	// 如果需要完整域名，可以在这里拼接 config.GetString("app.url")
	finalURL := localPath

	// 5. 更新业务表
	taskModel.UpdateProgress(90)

	// 更新角色表 (Character)
	if p.CharacterID > 0 {
		err = database.DB.Model(&characters.Characters{}).
			Where("id = ?", p.CharacterID).
			Update("avatar_url", finalURL).Error
	}

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("db update failed: %v", err))
		return err
	}

	// 6. 任务完成
	taskModel.MarkAsSuccess(fmt.Sprintf(`{"url": "%s"}`, finalURL))
	console.Success(fmt.Sprintf("任务[%d] - 图片生成并保存完成: %s", p.AsyncTaskID, finalURL))
	return nil
}
