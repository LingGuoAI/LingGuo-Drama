package scenes

import (
	"spiritFruit/app/models"
	"spiritFruit/app/models/projects"
	"spiritFruit/pkg/database"
)

// Scenes 结构体 场景/背景库
type Scenes struct {
	models.BaseModel
	ProjectId    *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;comment:所属项目ID;"`
	Name         *string `json:"name" form:"name" gorm:"column:name;comment:场景名称(地点+时间);size:100;"`
	Location     *string `json:"location" form:"location" gorm:"column:location;comment:地点;size:100;"`
	Time         *string `json:"time" form:"time" gorm:"column:time;comment:时间(日/夜);size:50;"`
	Atmosphere   *string `json:"atmosphere" form:"atmosphere" gorm:"column:atmosphere;comment:氛围描述;"`
	VisualPrompt *string `json:"visualPrompt" form:"visualPrompt" gorm:"column:visual_prompt;comment:AI绘画Prompt(纯背景);type:text;"`
	Status       *int8   `json:"status" form:"status" gorm:"default:1;column:status;comment:状态 1-待生成 2-生成中 3-已完成;"`

	// 关联关系
	Projectss *projects.Projects `json:"projects,omitempty" gorm:"foreignKey:ProjectId;references:ID"`
	models.CommonTimestampsField
}

func (Scenes) TableName() string {
	return "scenes"
}

// Create 创建场景项目
func (scenes *Scenes) Create() {
	database.DB.Create(&scenes)
}

// Save 保存场景项目
func (scenes *Scenes) Save() (rowsAffected int64) {
	result := database.DB.Save(&scenes)
	return result.RowsAffected
}

// Delete 删除场景项目
func (scenes *Scenes) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&scenes)
	return result.RowsAffected
}
