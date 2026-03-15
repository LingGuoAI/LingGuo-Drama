package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"spiritFruit/app/services"
	"strings"

	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/shot_frame_image" // 🔴 引入您提供的新模型
	"spiritFruit/app/models/shots"
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
	"spiritFruit/pkg/upload"
)

// HandleGenerateFrameImageTask 处理分镜帧图片生成
func HandleGenerateFrameImageTask(ctx context.Context, t *asynq.Task) error {
	var p myAsynq.GenerateFrameImagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 1. 获取并标记任务开始
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		return nil // 任务不存在
	}
	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始生成分镜图片 (ShotID: %d, Type: %s)", p.AsyncTaskID, p.ShotID, p.FrameType))

	// 2. 初始化 AI 配置
	taskModel.UpdateProgress(20)
	// ==========================================
	// 🔴 生图任务优先查库，后查 .env
	// ==========================================

	// 1) 先用 .env 环境变量进行基础兜底初始化
	aiConfig := openai.Config{
		Provider: config.GetString("ai.provider", "openai"),

		// OpenAI 配置
		OpenAIBaseURL: config.GetString("ai.openai.base_url"),
		OpenAIKey:     config.GetString("ai.openai.api_key"),
		OpenAIModel:   config.GetString("ai.openai.model"),

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
	errConfig, dbConfig := aiService.GetActiveConfigByType("image")

	if errConfig == nil && dbConfig.ID > 0 {
		providerName := *dbConfig.Provider
		baseURL := *dbConfig.BaseUrl
		apiKey := *dbConfig.ApiKey

		// 取 JSON 数组配置的第一个模型作为生图模型
		modelName := ""
		if len(dbConfig.Model) > 0 {
			modelName = dbConfig.Model[0]
		}

		// 动态覆盖兜底配置 (针对图片模型字段)
		switch providerName {
		case "openai", "getgoapi":
			aiConfig.Provider = "openai"
			aiConfig.OpenAIBaseURL = baseURL
			aiConfig.OpenAIKey = apiKey
			if modelName != "" {
				aiConfig.OpenAIModel = modelName
			}
		case "gemini", "google":
			aiConfig.Provider = "gemini"
			aiConfig.GeminiBaseURL = baseURL
			aiConfig.GeminiKey = apiKey
			if modelName != "" {
				aiConfig.GeminiModel = modelName
			}
		case "doubao", "volcengine":
			aiConfig.Provider = "doubao"
			aiConfig.DoubaoBaseURL = baseURL
			aiConfig.DoubaoKey = apiKey
			if modelName != "" {
				aiConfig.DoubaoImageModel = modelName // 赋值给图片模型字段
			}
		case "vertex":
			aiConfig.Provider = "vertex"
			aiConfig.VertexKey = apiKey
			if modelName != "" {
				aiConfig.VertexImageModel = modelName // 赋值给图片模型字段
			}
		default:
			aiConfig.Provider = "openai"
			aiConfig.OpenAIBaseURL = baseURL
			aiConfig.OpenAIKey = apiKey
			if modelName != "" {
				aiConfig.OpenAIModel = modelName
			}
		}

		console.Success(fmt.Sprintf("任务[%d] - 成功挂载数据库 AI 图片配置: Provider=%s, Model=%s", p.AsyncTaskID, providerName, modelName))
	} else {
		console.Warning(fmt.Sprintf("任务[%d] - 未命中数据库 AI 图片配置，将降级使用 .env 默认配置", p.AsyncTaskID))
	}
	aiProvider := openai.NewProvider(aiConfig)

	// 3. 调用 AI
	taskModel.UpdateProgress(40)
	// 组装最终 prompt
	imagePrompt := p.Prompt
	if !strings.Contains(imagePrompt, "cinematic lighting") {
		// 如果前端提示词没有这些高清词，稍微兜底补一下
		imagePrompt += ", cinematic lighting, highly detailed, realistic, 8k"
	}

	req := openai.ImageRequest{
		Prompt: imagePrompt,
		N:      1,
		Size:   "1024x1024", // 这里如果是首尾帧也许可以是宽屏 1024x1792
	}

	urls, err := aiProvider.GenerateImage(req)
	if err != nil {
		taskModel.MarkAsFailed(err)
		return err
	}
	if len(urls) == 0 {
		taskModel.MarkAsFailed(fmt.Errorf("no images generated"))
		return nil
	}

	// 4. 下载并保存到本地
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

	// 5. 将生成的图片写入 shot_frame_images 表
	taskModel.UpdateProgress(90)
	imagType := "shot" // 镜头图片
	newFrameImage := shot_frame_image.ShotFrameImages{
		ProjectId: &p.ProjectID,
		ShotId:    &p.ShotID,
		FrameType: &p.FrameType,
		ImageType: &imagType,
		ImageUrl:  &finalURL,
	}
	newFrameImage.Create()

	// 如果生成的是“首帧(first)”，自动将该图设置为镜头主表的封面 imageUrl
	if p.FrameType == "first" {
		database.DB.Model(&shots.Shots{}).Where("id = ?", p.ShotID).Update("image_url", finalURL)
	}

	// 6. 完成
	// 返回完整数据，前端可以直接解析拿去更新列表
	resultData := map[string]interface{}{
		"url":        finalURL,
		"shot_id":    p.ShotID,
		"frame_type": p.FrameType,
		"id":         newFrameImage.ID, // 数据库新插入记录的主键ID
	}
	resBytes, _ := json.Marshal(resultData)

	taskModel.MarkAsSuccess(string(resBytes))
	console.Success(fmt.Sprintf("任务[%d] - 分镜图片生成完成", p.AsyncTaskID))
	return nil
}
