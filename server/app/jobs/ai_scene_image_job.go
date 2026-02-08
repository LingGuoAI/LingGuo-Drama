package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/scenes" // 引入场景模型
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
	"spiritFruit/pkg/upload"
)

// HandleGenerateSceneImage 处理场景图片生成
func HandleGenerateSceneImage(ctx context.Context, t *asynq.Task) error {
	var p myAsynq.GenerateSceneImagePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 1. 获取并标记任务开始
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		return nil // 任务不存在
	}
	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始生成场景图片 (SceneID: %d)", p.AsyncTaskID, p.SceneID))

	// 2. 初始化 AI 配置
	taskModel.UpdateProgress(20)
	aiConfig := openai.Config{
		Provider:      config.GetString("ai.provider"),
		OpenAIBaseURL: config.GetString("ai.openai.base_url"),
		OpenAIKey:     config.GetString("ai.openai.api_key"),
		OpenAIModel:   "dall-e-3",
		GeminiBaseURL: config.GetString("ai.gemini.base_url"),
		GeminiKey:     config.GetString("ai.gemini.api_key"),
		GeminiModel:   "imagen-3.0-generate-001",
	}
	aiProvider := openai.NewProvider(aiConfig)

	// 3. 调用 AI
	taskModel.UpdateProgress(40)
	// 场景图通常是宽屏，如果有Dall-E-3支持，可改为 "1024x1792"
	req := openai.ImageRequest{
		Prompt: p.Prompt + ", cinematic lighting, highly detailed, realistic, 8k", // 增加一些画质词
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

	// 5. 更新 Scenes 表
	taskModel.UpdateProgress(90)
	// 根据前端逻辑，将图片URL存入 VisualPrompt 字段（或者您有独立的 Image 字段）
	err = database.DB.Model(&scenes.Scenes{}).
		Where("id = ?", p.SceneID).
		Update("visual_prompt", finalURL).Error

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("db update failed: %v", err))
		return err
	}

	// 6. 完成
	taskModel.MarkAsSuccess(fmt.Sprintf(`{"url": "%s"}`, finalURL))
	console.Success(fmt.Sprintf("任务[%d] - 场景图片生成完成", p.AsyncTaskID))
	return nil
}
