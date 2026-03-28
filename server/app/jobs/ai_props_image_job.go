package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"spiritFruit/app/services"
	"strings"

	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/props" // 引入道具模型
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
	"spiritFruit/pkg/upload"
)

// HandleGeneratePropImageTask 处理道具图片生成
func HandleGeneratePropImageTask(ctx context.Context, t *asynq.Task) error {
	// 1. 解析参数
	var p myAsynq.GeneratePropImagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 2. 获取并标记任务开始
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		console.Error(fmt.Sprintf("Task %d not found in DB", p.AsyncTaskID))
		return nil // 任务不存在，直接结束
	}
	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始生成道具图片 (PropID: %d)", p.AsyncTaskID, p.PropID))

	// 3. 初始化 AI 配置
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

	// 4. 构造 Prompt 并调用 AI
	taskModel.UpdateProgress(40)

	// 如果 Payload 中的 Prompt 为空，尝试从数据库获取兜底
	finalPrompt := p.Prompt
	if finalPrompt == "" {
		var prop props.Props
		if err := database.DB.First(&prop, p.PropID).Error; err == nil {
			if prop.ImagePrompt != nil && *prop.ImagePrompt != "" {
				finalPrompt = *prop.ImagePrompt
			} else {
				// 如果还没有提示词，用名称和描述拼接
				name := ""
				if prop.Name != nil {
					name = *prop.Name
				}
				desc := ""
				if prop.Description != nil {
					desc = *prop.Description
				}
				finalPrompt = fmt.Sprintf("A prop for a movie: %s. %s", name, desc)
			}
		}
	}

	// 道具图通常需要抠图或纯色背景，加上 white background 或 product shot 关键词
	enhancedPrompt := finalPrompt + ", product photography, white background, high quality, realistic, 8k, detailed texture"

	req := openai.ImageRequest{
		Prompt: enhancedPrompt,
		N:      1,
		Size:   "1024x1024", // 道具通常方形即可
	}

	console.Success(fmt.Sprintf("任务[%d] - Sending prompt: %s", p.AsyncTaskID, enhancedPrompt))

	urls, err := aiProvider.GenerateImage(req)
	if err != nil {
		taskModel.MarkAsFailed(err)
		return err
	}
	if len(urls) == 0 {
		taskModel.MarkAsFailed(fmt.Errorf("no images generated"))
		return nil
	}

	// 5. 下载并保存到本地
	taskModel.UpdateProgress(70)
	rawImageURL := urls[0]
	var localPath string
	var saveErr error

	if strings.HasPrefix(rawImageURL, "data:image") {
		localPath, saveErr = upload.SaveBase64Image(rawImageURL)
	} else {
		localPath, saveErr = upload.DownloadAndSave(rawImageURL)
	}

	if saveErr != nil {
		taskModel.MarkAsFailed(fmt.Errorf("save image failed: %v", saveErr))
		return saveErr
	}

	finalURL := localPath // 相对路径

	// 6. 更新 Props 表
	taskModel.UpdateProgress(90)

	// 更新 ImageUrl 字段，顺便回填 ImagePrompt (如果原本是空的)
	updates := map[string]interface{}{
		"image_url": finalURL,
	}
	// 如果使用了我们在代码里生成的兜底 Prompt，也可以选择存入数据库
	if p.Prompt == "" && finalPrompt != "" {
		updates["image_prompt"] = finalPrompt
	}

	err = database.DB.Model(&props.Props{}).
		Where("id = ?", p.PropID).
		Updates(updates).Error

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("db update failed: %v", err))
		return err
	}

	// 7. 完成
	taskModel.MarkAsSuccess(fmt.Sprintf(`{"url": "%s"}`, finalURL))
	console.Success(fmt.Sprintf("任务[%d] - 道具图片生成完成", p.AsyncTaskID))
	return nil
}
