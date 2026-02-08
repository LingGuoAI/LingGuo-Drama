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

// CreateImageGenerationTask 创建图片生成任务 (角色生图)
func (s *TaskService) CreateImageGenerationTask(projectID, charID uint64, prompt string) (*async_tasks.AsyncTask, error) {
	// 1. 构造 Payload
	payload := asynq.GenerateImagePayload{
		ProjectID:   projectID,
		CharacterID: charID,
		Prompt:      prompt,
	}
	payloadBytes, _ := json.Marshal(payload)

	// 2. 创建数据库记录
	task := async_tasks.AsyncTask{
		ProjectID: projectID,
		RelID:     charID, // 关联的角色ID
		Type:      asynq.TypeGenerateImage,
		Status:    async_tasks.StatusPending,
		Payload:   string(payloadBytes),
	}
	task.Create()

	// 3. 注入 ID
	payload.AsyncTaskID = task.ID

	// 4. 投递到 Asynq
	_, err := asynq.EnqueueGenerateImage(payload)
	if err != nil {
		task.MarkAsFailed(err)
		return &task, err
	}

	return &task, nil
}

// CreateGenerateCharactersTask 创建角色生成任务
func (s *TaskService) CreateGenerateCharactersTask(dramaID uint64, count int, outline string) (*async_tasks.AsyncTask, error) {
	// 1. 构造 Payload 数据
	payload := asynq.GenerateCharactersPayload{
		DramaID: dramaID,
		Count:   count,
		Outline: outline,
	}
	payloadBytes, _ := json.Marshal(payload)

	// 2. 先在数据库创建记录 (状态: Pending)
	task := async_tasks.AsyncTask{
		ProjectID: dramaID,                      // 关联项目ID
		Type:      asynq.TypeGenerateCharacters, // 使用 asynq 包中定义的类型常量
		Status:    async_tasks.StatusPending,
		Payload:   string(payloadBytes),
	}
	task.Create()

	// 3. 将数据库 ID 注入 Payload
	payload.AsyncTaskID = task.ID

	// 4. 投递到 Asynq
	_, err := asynq.EnqueueGenerateCharacters(payload)
	if err != nil {
		task.MarkAsFailed(err)
		return &task, err
	}

	return &task, nil
}

// CreateExtractScenesTask 创建场景提取任务
func (s *TaskService) CreateExtractScenesTask(projectID, episodeID uint64) (*async_tasks.AsyncTask, error) {
	// 1. 构造 Payload 数据
	payload := asynq.ExtractScenesPayload{
		EpisodeID: episodeID,
	}
	payloadBytes, _ := json.Marshal(payload)

	// 2. 先在数据库创建记录
	task := async_tasks.AsyncTask{
		ProjectID: projectID, // 尽量关联到项目ID，方便前端查询
		RelID:     episodeID, // 关联的章节ID
		Type:      asynq.TypeExtractScenes,
		Status:    async_tasks.StatusPending,
		Payload:   string(payloadBytes),
	}
	task.Create()

	// 3. 注入 ID
	payload.AsyncTaskID = task.ID

	// 4. 投递
	_, err := asynq.EnqueueExtractScenes(payload)
	if err != nil {
		task.MarkAsFailed(err)
		return &task, err
	}

	return &task, nil
}

// CreateSceneImageGenerationTask 创建场景图片生成任务
func (s *TaskService) CreateSceneImageGenerationTask(projectID, sceneID uint64, prompt string) (*async_tasks.AsyncTask, error) {
	// 1. 构造 Payload
	payload := asynq.GenerateSceneImagePayload{
		ProjectID: projectID,
		SceneID:   sceneID,
		Prompt:    prompt,
	}
	payloadBytes, _ := json.Marshal(payload)

	// 2. 创建数据库记录 (用于前端轮询)
	task := async_tasks.AsyncTask{
		ProjectID: projectID,
		RelID:     sceneID, // 关联的场景ID
		Type:      asynq.TypeGenerateSceneImage,
		Status:    async_tasks.StatusPending,
		Payload:   string(payloadBytes),
	}
	task.Create()

	// 3. 将数据库 ID 注入 Payload
	payload.AsyncTaskID = task.ID

	// 4. 投递到 Asynq
	_, err := asynq.EnqueueGenerateSceneImage(payload)
	if err != nil {
		task.MarkAsFailed(err)
		return &task, err
	}

	return &task, nil
}

// CreateGenerateShotsTask 创建分镜拆分任务
func (s *TaskService) CreateGenerateShotsTask(projectID, scriptID uint64, model string) (*async_tasks.AsyncTask, error) {
	// 1. 构造 Payload
	payload := asynq.GenerateShotsPayload{
		ProjectID: projectID,
		ScriptID:  scriptID,
		Model:     model,
	}
	payloadBytes, _ := json.Marshal(payload)

	// 2. 创建数据库记录
	task := async_tasks.AsyncTask{
		ProjectID: projectID,
		RelID:     scriptID, // 关联的剧本/分集ID
		Type:      asynq.TypeGenerateShots,
		Status:    async_tasks.StatusPending,
		Payload:   string(payloadBytes),
	}
	task.Create()

	// 3. 注入 Task ID
	payload.AsyncTaskID = task.ID

	// 4. 投递到 Asynq
	_, err := asynq.EnqueueGenerateShots(payload)
	if err != nil {
		task.MarkAsFailed(err)
		return &task, err
	}

	return &task, nil
}
