package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/scenes"
	"spiritFruit/app/models/scripts"
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
	"spiritFruit/pkg/prompt"
	"spiritFruit/pkg/utils"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// HandleExtractScenes 处理场景提取任务
func HandleExtractScenes(ctx context.Context, t *asynq.Task) error {
	// 1. 解析参数
	var p myAsynq.ExtractScenesPayload
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
	console.Success(fmt.Sprintf("任务[%d] - 开始提取场景", p.AsyncTaskID))

	// 3. 获取剧本章节内容
	var episode scripts.Scripts
	// Preload Projectss 拿到风格 Style
	if err := database.DB.Preload("Projectss").First(&episode, p.ScriptID).Error; err != nil {
		err = fmt.Errorf("episode script not found: %v", err)
		taskModel.MarkAsFailed(err)
		return nil
	}

	if episode.Content == nil || *episode.Content == "" {
		err := fmt.Errorf("script content is empty")
		taskModel.MarkAsFailed(err)
		return nil
	}

	// 获取风格
	dramaStyle := "realistic"
	if episode.Projectss != nil && episode.Projectss.Style != nil {
		dramaStyle = *episode.Projectss.Style
	}

	// [Stage 2] 准备 Prompt 和 AI 配置，进度 -> 20%
	taskModel.UpdateProgress(20)

	promptGen := prompt.NewGenerator()
	systemPrompt := promptGen.GetSceneExtractionPrompt(dramaStyle)
	contentLabel := "【剧本内容】"
	formatInstructions := getSceneFormatInstructions(promptGen.IsEnglish())
	finalPrompt := fmt.Sprintf("%s\n\n%s\n%s\n\n%s", systemPrompt, contentLabel, *episode.Content, formatInstructions)

	// AI 配置
	aiConfig := openai.Config{
		Provider:      config.GetString("ai.provider"),
		OpenAIBaseURL: config.GetString("ai.openai.base_url"),
		OpenAIKey:     config.GetString("ai.openai.api_key"),
		OpenAIModel:   config.GetString("ai.openai.model"),
		GeminiBaseURL: config.GetString("ai.gemini.base_url"),
		GeminiKey:     config.GetString("ai.gemini.api_key"),
		GeminiModel:   config.GetString("ai.gemini.model"),
	}
	aiProvider := openai.NewProvider(aiConfig)

	// [Stage 3] 发起 AI 请求，进度 -> 30%
	taskModel.UpdateProgress(30)
	console.Success(fmt.Sprintf("任务[%d] - Sending prompt to AI", p.AsyncTaskID))

	req := openai.ScriptRequest{
		Messages: []openai.ChatMessage{
			// 对于场景提取，通常只需要一次对话即可，这里为了简单直接发一条 User 消息
			// 也可以拆分为 System 和 User
			{Role: "user", Content: finalPrompt},
		},
		Temperature: 0.7,
	}

	aiResp, err := aiProvider.GenerateScript(req)
	if err != nil {
		console.Error(fmt.Sprintf("AI提取失败: %v", err))
		taskModel.MarkAsFailed(err)
		return err
	}

	// [Stage 4] 解析结果，进度 -> 60%
	taskModel.UpdateProgress(60)

	type BackgroundInfo struct {
		Location   string `json:"location"`
		Time       string `json:"time"`
		Atmosphere string `json:"atmosphere"`
		Prompt     string `json:"prompt"`
	}

	var backgrounds []BackgroundInfo

	// 兼容 Array 和 Object 两种返回格式
	if err := utils.SafeParseAIJSON(aiResp, &backgrounds); err != nil {
		var wrapper struct {
			Backgrounds []BackgroundInfo `json:"backgrounds"`
		}
		if err2 := utils.SafeParseAIJSON(aiResp, &wrapper); err2 != nil {
			err = fmt.Errorf("failed to parse scene JSON: %v", err)
			taskModel.MarkAsFailed(err)
			return nil
		}
		backgrounds = wrapper.Backgrounds
	}

	// [Stage 5] 数据入库，进度 -> 80%
	taskModel.UpdateProgress(80)

	count := 0
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		for _, bg := range backgrounds {
			// 构造场景名
			sceneName := fmt.Sprintf("%s-%s", bg.Location, bg.Time)

			var existCount int64
			// ProjectId 是指针，需要解引用或者直接用
			if episode.ProjectId == nil {
				return fmt.Errorf("project id is nil")
			}
			projID := *episode.ProjectId

			tx.Model(&scenes.Scenes{}).Where("project_id = ? AND name = ?", projID, sceneName).Count(&existCount)

			if existCount == 0 {
				loc := bg.Location
				tm := bg.Time
				atm := bg.Atmosphere
				prt := bg.Prompt
				status := int8(1) // 1-待生成

				newScene := scenes.Scenes{
					ProjectId:    &projID,
					Name:         &sceneName,
					Location:     &loc,
					Time:         &tm,
					Atmosphere:   &atm,
					VisualPrompt: &prt,
					Status:       &status,
				}
				if err := tx.Create(&newScene).Error; err != nil {
					return err
				}
				count++
			}
		}
		return nil
	})

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("database transaction failed: %v", err))
		return err
	}

	// [Stage 6] 全部完成，进度 -> 100%
	resultInfo := fmt.Sprintf(`{"extracted_count": %d, "provider": "%s"}`, count, aiConfig.Provider)
	taskModel.MarkAsSuccess(resultInfo)

	console.Success(fmt.Sprintf("任务[%d] - 场景提取完成，新增 %d 个场景", p.AsyncTaskID, count))
	return nil
}

// 辅助函数保持不变
func getSceneFormatInstructions(isEnglish bool) string {
	if isEnglish {
		return `[Output JSON Format]
{
  "backgrounds": [
    {
      "location": "Location Name",
      "time": "Time Description",
      "atmosphere": "Atmosphere",
      "prompt": "Detailed English image generation prompt..."
    }
  ]
}`
	}
	return `【输出JSON格式】
{
  "backgrounds": [
    {
      "location": "地点名称",
      "time": "时间描述",
      "atmosphere": "氛围",
      "prompt": "详细的中文图片生成提示词，纯背景，无人物..."
    }
  ]
}`
}
