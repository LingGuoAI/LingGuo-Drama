package asynq

const (
	TypeGenerateScript     = "drama:generate_script"
	TypeGenerateImage      = "drama:generate_image"
	TypeGenerateCharacters = "ai:generate:characters" // 提取角色
	TypeExtractScenes      = "ai:extract:scenes"      // 提取场景
	TypeGenerateSceneImage = "generate:scene:image"
	TypeGenerateShots      = "generate:shots"
	TypeExtractProps       = "ai:extract_props"
	TypeGeneratePropImage  = "ai:generate_prop_image"
)

// GenerateScriptPayload
type GenerateScriptPayload struct {
	AsyncTaskID uint64 `json:"asyncTaskId"` // 关键：关联的任务表ID
	ProjectID   uint64 `json:"projectId"`
	ScriptID    uint64 `json:"scriptId"` // 具体要更新哪个剧本
	Prompt      string `json:"prompt"`
}

// GenerateImagePayload 生成图片的参数
type GenerateImagePayload struct {
	AsyncTaskID uint64 `json:"asyncTaskId"` // 关联的任务ID
	ProjectID   uint64 `json:"projectId"`   // 关联的项目ID

	// 互斥参数：要么传 CharacterID (角色定妆)，要么传 ShotID (分镜图)
	CharacterID uint64 `json:"characterId,omitempty"`
	ShotID      uint64 `json:"shotId,omitempty"`

	Prompt string `json:"prompt"`
}

type GenerateCharactersPayload struct {
	AsyncTaskID uint64 `json:"asyncTaskId"`
	ProjectID   uint64 `json:"projectId"`
	Count       int    `json:"count"`
	Outline     string `json:"outline"`
}

type ExtractScenesPayload struct {
	AsyncTaskID uint64 `json:"asyncTaskId"`
	ScriptID    uint64 `json:"scriptId"`
}

// GenerateSceneImagePayload 场景图片生成载荷
type GenerateSceneImagePayload struct {
	AsyncTaskID uint64 `json:"asyncTaskId"`
	ProjectID   uint64 `json:"projectId"`
	SceneID     uint64 `json:"sceneId"`
	Prompt      string `json:"prompt"`
}

// GenerateShotsPayload 分镜生成载荷
type GenerateShotsPayload struct {
	AsyncTaskID uint64 `json:"asyncTaskId"`
	ProjectID   uint64 `json:"projectId"`
	ScriptID    uint64 `json:"scriptId"`
	Model       string `json:"model"` // 可选：指定模型
}

// GeneratePropImagePayload 道具生图载荷
type GeneratePropImagePayload struct {
	AsyncTaskID uint64 `json:"async_task_id"`
	ProjectID   uint64 `json:"project_id"`
	PropID      uint64 `json:"prop_id"`
	Prompt      string `json:"prompt"`
}

// ExtractPropsPayload 提取道具载荷
type ExtractPropsPayload struct {
	AsyncTaskID uint64 `json:"async_task_id"`
	ProjectID   uint64 `json:"project_id"`
	EpisodeID   uint64 `json:"episode_id"` // 根据哪一集剧本提取
}
