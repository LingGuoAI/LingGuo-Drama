package source

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

type Source struct {
	models.BaseModel
	ProjectId  *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;not null;default:0;comment:项目id;"`
	ScriptId   *uint64 `json:"scriptId" form:"scriptId" gorm:"column:script_id;default:0;comment:集数id"`
	ShotId     *uint64 `json:"shotId" form:"shotId" gorm:"column:shot_id;not null;default:0;comment:分镜ID;"`
	ShotNumber *uint64 `json:"shotNumber" form:"shotNumber" gorm:"column:shot_number;comment:第几个镜头"`
	VideoUrl   *string `json:"videoUrl" form:"videoUrl" gorm:"column:video_url;comment:视频路径"`
	models.CommonTimestampsField
}

// TableName 素材表
func (Source) TableName() string {
	return "sources"
}

// Create 创建素材
func (source *Source) Create() {
	database.DB.Create(&source)
}

// Save 保存素材
func (source *Source) Save() (rowsAffected int64) {
	result := database.DB.Save(&source)
	return result.RowsAffected
}

// Delete 删除素材
func (source *Source) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&source)
	return result.RowsAffected
}
