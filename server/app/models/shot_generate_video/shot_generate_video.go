package shot_generate_video

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

type ShotGenerateVideo struct {
	models.BaseModel
	ProjectId *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;not null;default:0;comment:项目id;"`
	ScriptId  *uint64 `json:"scriptId" form:"scriptId" gorm:"column:script_id;default:0;comment:集数id"`
	ShotId    *uint64 `json:"shotId" form:"shotId" gorm:"column:shot_id;not null;default:0;comment:分镜ID;"`
	VideoUrl  *string `json:"videoUrl" form:"videoUrl" gorm:"column:video_url;comment:视频"`
	models.CommonTimestampsField
}

// TableName 分镜生成视频表
func (ShotGenerateVideo) TableName() string {
	return "shot_generate_videos"
}

// Create 创建分镜生成视频
func (shotGenerateVideo *ShotGenerateVideo) Create() {
	database.DB.Create(&shotGenerateVideo)
}

// Save 保存分镜生成视频
func (shotGenerateVideo *ShotGenerateVideo) Save() (rowsAffected int64) {
	result := database.DB.Save(&shotGenerateVideo)
	return result.RowsAffected
}

// Delete 删除分镜生成视频
func (shotGenerateVideo *ShotGenerateVideo) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&shotGenerateVideo)
	return result.RowsAffected
}
