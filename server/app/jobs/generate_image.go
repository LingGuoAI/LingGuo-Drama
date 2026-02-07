package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/characters" // 或者是角色表
	"spiritFruit/app/models/shots"      // 假设我们要更新分镜表
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

	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		return nil
	}

	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始生成图片", p.AsyncTaskID))

	// 进度 20%
	taskModel.UpdateProgress(20)

	aiConfig := openai.Config{
		Provider:      config.GetString("ai.provider"),
		OpenAIBaseURL: config.GetString("ai.openai.base_url"),
		OpenAIKey:     config.GetString("ai.openai.api_key"),
		OpenAIModel:   "dall-e-3", // 生图通常指定模型

		GeminiBaseURL: config.GetString("ai.gemini.base_url"),
		GeminiKey:     config.GetString("ai.gemini.api_key"),
		GeminiModel:   "imagen-3.0-generate-001",
	}

	aiProvider := openai.NewProvider(aiConfig)

	// 进度 40% - 请求 AI
	taskModel.UpdateProgress(40)

	req := openai.ImageRequest{
		Prompt: p.Prompt,
		N:      1,
		Size:   "1024x1024",
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

	finalImageURL := urls[0]

	// TODO: 如果使用的是 Gemini 且返回的是 Base64，建议在此处调用 UploadService 上传到 OSS
	// if strings.HasPrefix(finalImageURL, "data:image") {
	//     finalImageURL, err = services.UploadBase64Image(finalImageURL)
	//     if err != nil { ... }
	// }

	// 进度 80% - 更新业务表
	taskModel.UpdateProgress(80)

	// 根据传入的 ID 类型更新不同的表
	if p.CharacterID > 0 {
		// 更新角色表
		err = database.DB.Model(&characters.Characters{}).
			Where("id = ?", p.CharacterID).
			Update("avatar_url", finalImageURL).Error
	} else if p.ShotID > 0 {
		// 更新分镜表
		err = database.DB.Model(&shots.Shots{}).
			Where("id = ?", p.ShotID).
			Update("image_url", finalImageURL).Error
	}

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("db update failed: %v", err))
		return err
	}

	taskModel.MarkAsSuccess(finalImageURL) // 结果里存一下 URL
	console.Success(fmt.Sprintf("任务[%d] - 图片生成完成", p.AsyncTaskID))
	return nil
}
