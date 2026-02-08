package shots

import (
	"spiritFruit/app/models"
	"spiritFruit/app/models/projects"
	"spiritFruit/app/models/scenes"
	"spiritFruit/app/models/scripts"
	"spiritFruit/pkg/database"
)

// Shots 结构体 镜头表
type Shots struct {
	models.BaseModel
	ProjectId *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;comment:所属项目ID;"`
	ScriptId  *uint64 `json:"scriptId" form:"scriptId" gorm:"column:script_id;comment:所属剧本ID;"`
	// 新增关联场景ID
	SceneId *uint64 `json:"sceneId" form:"sceneId" gorm:"column:scene_id;comment:关联场景ID;"`

	SequenceNo *uint64 `json:"sequenceNo" form:"sequenceNo" gorm:"default:0;column:sequence_no;comment:镜头序号;"`

	Title *string `json:"title" form:"title" gorm:"column:title;comment:镜头标题;size:100;"`

	ShotType       *string `json:"shotType" form:"shotType" gorm:"column:shot_type;comment:景别;size:50;"`
	CameraMovement *string `json:"cameraMovement" form:"cameraMovement" gorm:"column:camera_movement;comment:运镜;size:50;"`
	Angle          *string `json:"angle" form:"angle" gorm:"column:angle;comment:视角;size:50;"`

	Time     *string `json:"time" form:"time" gorm:"column:time;comment:具体时间描述;size:255;"`
	Location *string `json:"location" form:"location" gorm:"column:location;comment:具体地点描述;size:255;"`
	Action   *string `json:"action" form:"action" gorm:"column:action;comment:人物动作描述;type:text;"`

	Dialogue   *string `json:"dialogue" form:"dialogue" gorm:"column:dialogue;comment:台词;"`
	VisualDesc *string `json:"visualDesc" form:"visualDesc" gorm:"column:visual_desc;comment:画面结果(Visual Result);type:text;"`
	Atmosphere *string `json:"atmosphere" form:"atmosphere" gorm:"column:atmosphere;comment:氛围描述;"`

	// 生图/视频相关
	ImagePrompt *string `json:"imagePrompt" form:"imagePrompt" gorm:"column:image_prompt;comment:绘画Prompt;type:text;"`
	VideoPrompt *string `json:"videoPrompt" form:"videoPrompt" gorm:"column:video_prompt;comment:视频生成Prompt;type:text;"`
	AudioPrompt *string `json:"audioPrompt" form:"audioPrompt" gorm:"column:audio_prompt;comment:音效/BGM提示词;"`

	ImageUrl *string `json:"imageUrl" form:"imageUrl" gorm:"column:image_url;comment:分镜图;size:1024;"`
	VideoUrl *string `json:"videoUrl" form:"videoUrl" gorm:"column:video_url;comment:最终视频;size:1024;"`
	AudioUrl *string `json:"audioUrl" form:"audioUrl" gorm:"column:audio_url;comment:配音;size:1024;"`

	DurationMs *uint64 `json:"durationMs" form:"durationMs" gorm:"default:3000;column:duration_ms;"`
	Status     *int8   `json:"status" form:"status" gorm:"default:0;column:status;"`

	// 关联关系
	Projectss *projects.Projects `json:"projects,omitempty" gorm:"foreignKey:ProjectId;references:ID"`
	Scriptss  *scripts.Scripts   `json:"scripts,omitempty" gorm:"foreignKey:ScriptId;references:ID"`
	Sceness   *scenes.Scenes     `json:"scenes,omitempty" gorm:"foreignKey:SceneId;references:ID"` // 新增关联
	models.CommonTimestampsField
}

// TableName 镜头表 Shots自定义表名 shots
func (Shots) TableName() string {
	return "shots"
}

// Create 创建镜头表
func (shots *Shots) Create() {
	database.DB.Create(&shots)
}

// Save 保存镜头表
func (shots *Shots) Save() (rowsAffected int64) {
	result := database.DB.Save(&shots)
	return result.RowsAffected
}

// Delete 删除镜头表
func (shots *Shots) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&shots)
	return result.RowsAffected
}
