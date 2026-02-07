package services

import (
	"encoding/json"
	"spiritFruit/app/models/async_tasks"
	"spiritFruit/pkg/asynq"
)

type TaskService struct{}

// CreateScriptGenerationTask 创建剧本生成任务
func (s *TaskService) CreateScriptGenerationTask(projectID, scriptID uint64, prompt string) (*async_tasks.AsyncTask, error) {
	// 1. 构造 Payload 数据
	payload := asynq.GenerateScriptPayload{
		ProjectID: projectID,
		ScriptID:  scriptID,
		Prompt:    prompt,
	}
	payloadBytes, _ := json.Marshal(payload)

	// 2. 先在数据库创建记录 (状态: Pending)
	task := async_tasks.AsyncTask{
		ProjectID: projectID,
		RelID:     scriptID,
		Type:      async_tasks.TypeGenerateScript,
		Status:    async_tasks.StatusPending,
		Payload:   string(payloadBytes),
	}
	task.Create()

	// 3. 将数据库 ID 注入 Payload
	payload.AsyncTaskID = task.ID

	// 4. 投递到 Asynq (传入包含 TaskID 的 payload)
	_, err := asynq.EnqueueGenerateScript(payload)
	if err != nil {
		// 如果投递失败，标记任务为失败
		task.MarkAsFailed(err)
		return &task, err
	}

	return &task, nil
}

// CreateImageGenerationTask 创建图片生成任务 (类似)
func (s *TaskService) CreateImageGenerationTask(projectID, charID uint64, prompt string) (*async_tasks.AsyncTask, error) {
	// ... 逻辑同上，替换 Payload 类型即可
	return nil, nil // 省略实现
}
