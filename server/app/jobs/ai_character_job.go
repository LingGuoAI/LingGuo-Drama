package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/characters"
	"spiritFruit/app/models/projects"
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
	"spiritFruit/pkg/prompt"
	"spiritFruit/pkg/utils"

	"github.com/hibiken/asynq"
)

// HandleGenerateCharacters 处理角色生成任务
func HandleGenerateCharacters(ctx context.Context, t *asynq.Task) error {
	// 1. 解析参数
	var p myAsynq.GenerateCharactersPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 2. 获取任务并标记开始
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		console.Error(fmt.Sprintf("Task %d not found in DB", p.AsyncTaskID))
		return nil // 任务不存在，直接结束
	}

	// [Stage 1] 状态变更为 Processing，进度 -> 10%
	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始生成角色", p.AsyncTaskID))

	// 3. 获取项目信息
	var projectModel projects.Projects
	if err := database.DB.First(&projectModel, p.ProjectID).Error; err != nil {
		err = fmt.Errorf("project not found: %v", err)
		taskModel.MarkAsFailed(err)
		return nil // 项目不存在，无需重试
	}

	// 处理字段默认值
	projectTitle := ""
	if projectModel.Title != nil {
		projectTitle = *projectModel.Title
	}
	projectDesc := ""
	if projectModel.Description != nil {
		projectDesc = *projectModel.Description
	}
	projectGenre := "都市"
	if projectModel.Genre != nil {
		projectGenre = *projectModel.Genre
	}
	projectStyle := "realistic"
	if projectModel.Style != nil {
		projectStyle = *projectModel.Style
	}

	// [Stage 2] 准备 Prompt 和 AI 配置，进度 -> 20%
	taskModel.UpdateProgress(20)

	// 准备 Prompt
	promptGen := prompt.NewGenerator()
	systemPrompt := promptGen.GetCharacterExtractionPrompt(projectStyle)

	outlineText := p.Outline
	if outlineText == "" {
		outlineText = fmt.Sprintf("剧名：%s\n简介：%s\n类型：%s", projectTitle, projectDesc, projectGenre)
	}
	userPrompt := promptGen.FormatUserPrompt("character_request", outlineText, p.Count)

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
	if aiConfig.OpenAIModel == "" {
		aiConfig.OpenAIModel = "gpt-4-turbo"
	}

	aiProvider := openai.NewProvider(aiConfig)

	// [Stage 3] 发起 AI 请求，进度 -> 30%
	taskModel.UpdateProgress(30)
	console.Success(fmt.Sprintf("任务[%d] - Sending prompt to AI", p.AsyncTaskID))

	// 构造请求
	req := openai.ScriptRequest{
		Messages: []openai.ChatMessage{
			{Role: "system", Content: systemPrompt}, // 注意：部分模型(如Gemini)可能自动合并system到user
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
	}

	// 调用 AI
	aiResp, err := aiProvider.GenerateScript(req)
	if err != nil {
		console.Error(fmt.Sprintf("AI生成失败: %v", err))
		taskModel.MarkAsFailed(err)
		return err // 返回 err 触发重试
	}
	// [Stage 4] 解析结果，进度 -> 60%
	taskModel.UpdateProgress(60)

	var aiResult []struct {
		Name        string `json:"name"`
		Role        string `json:"role"`
		Description string `json:"description"`
		Personality string `json:"personality"`
		Appearance  string `json:"appearance"`
	}

	if err := utils.SafeParseAIJSON(aiResp, &aiResult); err != nil {
		err = fmt.Errorf("failed to parse AI response: %v", err)
		taskModel.MarkAsFailed(err)
		return nil // 解析失败通常是 AI 输出格式错，重试意义不大
	}

	// [Stage 5] 数据入库，进度 -> 80%
	taskModel.UpdateProgress(80)

	tx := database.DB.Begin()
	count := 0
	projID := p.ProjectID

	for _, char := range aiResult {
		// 查重
		var existCount int64
		tx.Model(&characters.Characters{}).Where("project_id = ? AND name = ?", projID, char.Name).Count(&existCount)
		if existCount > 0 {
			continue
		}

		roleType := char.Role
		pers := char.Personality
		appDesc := char.Appearance
		// 简单的 Visual Prompt 生成逻辑
		visualPrompt := fmt.Sprintf("%s, %s, %s", char.Appearance, projectStyle, "high quality, best quality")

		newChar := characters.Characters{
			ProjectId:      &projID,
			Name:           &char.Name,
			RoleType:       &roleType,
			Personality:    &pers,
			AppearanceDesc: &appDesc,
			VisualPrompt:   &visualPrompt,
		}

		if err := tx.Create(&newChar).Error; err != nil {
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

	console.Success(fmt.Sprintf("任务[%d] - 角色生成完成，新增 %d 个角色", p.AsyncTaskID, count))
	return nil
}
