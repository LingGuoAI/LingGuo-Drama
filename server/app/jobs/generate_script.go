package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/scripts"
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
)

func HandleGenerateScript(ctx context.Context, t *asynq.Task) error {
	// 1. 解析参数
	var p myAsynq.GenerateScriptPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 2. 获取任务并标记开始
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		// 如果数据库没找到任务，直接返回 nil 结束任务，避免无限重试
		console.Error(fmt.Sprintf("Task %d not found in DB", p.AsyncTaskID))
		return nil
	}

	// [Stage 1] 状态变更为 Processing，进度 -> 10%
	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始生成剧本", p.AsyncTaskID))

	// 3. 准备配置并初始化 AI 客户端
	// [Stage 2] 准备配置，进度 -> 20%
	taskModel.UpdateProgress(20)

	aiConfig := openai.Config{
		// 从配置文件读取，例如 config/app.go 或 .env
		Provider: config.GetString("ai.provider"), // "openai" 或 "gemini"

		OpenAIBaseURL: config.GetString("ai.openai.base_url"),
		OpenAIKey:     config.GetString("ai.openai.api_key"),
		OpenAIModel:   config.GetString("ai.openai.model"),

		GeminiBaseURL: config.GetString("ai.gemini.base_url"),
		GeminiKey:     config.GetString("ai.gemini.api_key"),
		GeminiModel:   config.GetString("ai.gemini.model"),
	}

	// 使用工厂方法创建 Provider
	aiProvider := openai.NewProvider(aiConfig)

	// [Stage 3] 发起 AI 请求
	taskModel.UpdateProgress(30)

	// 构造请求参数
	req := openai.ScriptRequest{
		Messages: []openai.ChatMessage{
			{Role: "user", Content: p.Prompt},
		},
		Temperature: 0.7, // 设置创意度
	}

	// 调用接口
	content, err := aiProvider.GenerateScript(req)
	if err != nil {
		// 记录失败原因
		console.Error(fmt.Sprintf("AI生成失败: %v", err))
		taskModel.MarkAsFailed(err)
		return err // 返回 err 触发 Asynq 重试
	}

	// [Stage 4] AI 生成完毕，准备写入数据库，进度 -> 80%
	taskModel.UpdateProgress(80)

	// 4. 业务数据落库 (更新 scripts 表)
	// 这里只更新内容，如果需要更新标题等可以解析 content
	err = database.DB.Model(&scripts.Scripts{}).
		Where("id = ?", p.ScriptID).
		Updates(map[string]interface{}{
			"content":   content,
			"is_locked": 0, // 确保未锁定，允许用户修改
		}).Error

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("save script content failed: %v", err))
		return err
	}

	// [Stage 5] 全部完成，进度 -> 100%
	// 可以在 Result 中存储一些元数据
	resultInfo := fmt.Sprintf(`{"length": %d, "provider": "%s"}`, len(content), aiConfig.Provider)
	taskModel.MarkAsSuccess(resultInfo)

	console.Success(fmt.Sprintf("任务[%d] - 剧本生成完成", p.AsyncTaskID))
	return nil
}
