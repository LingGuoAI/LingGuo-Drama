package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/projects" // 如果需要项目风格
	"spiritFruit/app/models/props"
	"spiritFruit/app/models/scripts"
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
	"spiritFruit/pkg/prompt"
	"spiritFruit/pkg/utils"

	"github.com/hibiken/asynq"
)

// HandleExtractPropsTask 处理道具提取任务
func HandleExtractPropsTask(ctx context.Context, t *asynq.Task) error {
	// 1. 解析参数
	var p myAsynq.ExtractPropsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 2. 获取任务并标记开始
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		console.Error(fmt.Sprintf("Task %d not found in DB", p.AsyncTaskID))
		return nil
	}

	// [Stage 1] 状态变更为 Processing，进度 -> 10%
	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始从剧本提取道具", p.AsyncTaskID))

	// 3. 获取剧本内容
	var scriptModel scripts.Scripts
	if err := database.DB.First(&scriptModel, p.EpisodeID).Error; err != nil {
		err = fmt.Errorf("script not found: %v", err)
		taskModel.MarkAsFailed(err)
		return nil
	}

	// 尝试获取项目风格 (用于 Prompt 优化，可选)
	var projectModel projects.Projects
	projectStyle := "realistic" // 默认风格
	if scriptModel.ProjectId != nil {
		if err := database.DB.First(&projectModel, *scriptModel.ProjectId).Error; err == nil {
			if projectModel.Style != nil {
				projectStyle = *projectModel.Style
			}
		}
	}

	// [Stage 2] 准备 Prompt 和 AI 配置，进度 -> 20%
	taskModel.UpdateProgress(20)

	promptGen := prompt.NewGenerator()
	systemPrompt := promptGen.GetPropExtractionPrompt(projectStyle)

	userPrompt := fmt.Sprintf("剧本内容：\n%s", *scriptModel.Content)

	// 准备 AI 配置
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
	}
	aiProvider := openai.NewProvider(aiConfig)

	// [Stage 3] 发起 AI 请求，进度 -> 30%
	taskModel.UpdateProgress(30)
	console.Success(fmt.Sprintf("任务[%d] - Sending prompt to AI", p.AsyncTaskID))

	req := openai.ScriptRequest{
		Messages: []openai.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.5, // 提取任务温度低一点更稳定
	}

	aiResp, err := aiProvider.GenerateScript(req)
	if err != nil {
		console.Error(fmt.Sprintf("AI生成失败: %v", err))
		taskModel.MarkAsFailed(err)
		return err
	}
	// [Stage 4] 解析结果，进度 -> 60%
	taskModel.UpdateProgress(60)

	type ExtractedProp struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		ImagePrompt string `json:"image_prompt"`
	}
	var aiResult []ExtractedProp

	if err := utils.SafeParseAIJSON(aiResp, &aiResult); err != nil {
		err = fmt.Errorf("failed to parse AI response: %v", err)
		taskModel.MarkAsFailed(err)
		return nil
	}

	// [Stage 5] 数据入库，进度 -> 80%
	taskModel.UpdateProgress(80)

	tx := database.DB.Begin()
	count := 0
	projID := p.ProjectID

	for _, item := range aiResult {
		// 查重
		var existCount int64
		tx.Model(&props.Props{}).Where("project_id = ? AND name = ?", projID, item.Name).Count(&existCount)
		if existCount > 0 {
			continue
		}

		newProp := props.Props{
			ProjectId:   &projID,
			Name:        &item.Name,
			Type:        &item.Type,
			Description: &item.Description,
			ImagePrompt: &item.ImagePrompt,
		}

		if err := tx.Create(&newProp).Error; err != nil {
			tx.Rollback()
			err = fmt.Errorf("db create failed: %v", err)
			taskModel.MarkAsFailed(err)
			return err
		}
		count++
	}
	tx.Commit()

	// [Stage 6] 全部完成，进度 -> 100%
	resultInfo := fmt.Sprintf(`{"generated_count": %d, "provider": "%s"}`, count, aiConfig.Provider)
	taskModel.MarkAsSuccess(resultInfo)

	console.Success(fmt.Sprintf("任务[%d] - 道具提取完成，新增 %d 个道具", p.AsyncTaskID, count))
	return nil
}
