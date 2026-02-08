package asynq

const (
	TypeGenerateScript     = "drama:generate_script"
	TypeGenerateImage      = "drama:generate_image"
	TypeGenerateCharacters = "ai:generate:characters" // 提取角色
	TypeExtractScenes      = "ai:extract:scenes"      // 提取场景
	TypeGenerateSceneImage = "generate:scene:image"
	TypeGenerateShots      = "generate:shots"
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

type GenerateCharactersPayload struct {
	AsyncTaskID uint64 `json:"async_task_id"`
	DramaID     uint64 `json:"drama_id"`
	Count       int    `json:"count"`
	Outline     string `json:"outline"`
}

type ExtractScenesPayload struct {
	AsyncTaskID uint64 `json:"async_task_id"`
	EpisodeID   uint64 `json:"episode_id"`
}

// GenerateSceneImagePayload 场景图片生成载荷
type GenerateSceneImagePayload struct {
	AsyncTaskID uint64 `json:"async_task_id"`
	ProjectID   uint64 `json:"project_id"`
	SceneID     uint64 `json:"scene_id"`
	Prompt      string `json:"prompt"`
}

// GenerateShotsPayload 分镜生成载荷
type GenerateShotsPayload struct {
	AsyncTaskID uint64 `json:"async_task_id"`
	ProjectID   uint64 `json:"project_id"`
	ScriptID    uint64 `json:"script_id"` // 对应 episode/剧本 ID
	Model       string `json:"model"`     // 可选：指定模型
}
