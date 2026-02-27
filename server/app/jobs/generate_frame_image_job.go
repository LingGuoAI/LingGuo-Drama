package jobs

import (
	"context"
	"encoding/json"
	"fmt"
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

	newFrameImage := shot_frame_image.ShotFrameImages{
		ProjectId: &p.ProjectID,
		ShotId:    &p.ShotID,
		FrameType: &p.FrameType,
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
