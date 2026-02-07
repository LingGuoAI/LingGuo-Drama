package asynq

const (
	TypeGenerateScript = "drama:generate_script"
	TypeGenerateImage  = "drama:generate_image"
)

// GenerateScriptPayload
type GenerateScriptPayload struct {
	AsyncTaskID uint64 `json:"async_task_id"` // 关键：关联的任务表ID
	ProjectID   uint64 `json:"project_id"`
	ScriptID    uint64 `json:"script_id"` // 具体要更新哪个剧本
	Prompt      string `json:"prompt"`
}

// GenerateImagePayload 生成图片的参数
type GenerateImagePayload struct {
	AsyncTaskID uint64 `json:"async_task_id"` // 关联的任务ID
	ProjectID   uint64 `json:"project_id"`    // 关联的项目ID

	// 互斥参数：要么传 CharacterID (角色定妆)，要么传 ShotID (分镜图)
	CharacterID uint64 `json:"character_id,omitempty"`
	ShotID      uint64 `json:"shot_id,omitempty"`

	Prompt string `json:"prompt"`
}
