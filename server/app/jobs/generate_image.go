package jobs

import (
	"context"
	"encoding/json"
	"fmt"
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
	// 准备 AI 配置
	aiConfig := openai.Config{
		Provider: config.GetString("ai.provider", "openai"), // 提供默认值

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
