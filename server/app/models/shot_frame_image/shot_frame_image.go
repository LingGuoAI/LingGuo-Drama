package shot_frame_image

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

type ShotFrameImages struct {
	models.BaseModel
	ProjectId *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;not null;default:0;comment:项目id;"`
	ShotId    *uint64 `json:"shotId" form:"shotId" gorm:"column:shot_id;not null;default:0;comment:分镜ID;"`
	FrameType *string `json:"frameType" form:"frameType" gorm:"column:frame_type;comment:帧类型"` // first-首帧 last-尾帧 action-动作序列 key-关键帧
	ImageUrl  *string `json:"imageUrl" form:"imageUrl" gorm:"column:image_url;comment:图片"`
	models.CommonTimestampsField
}

// TableName 分镜图片表
func (ShotFrameImages) TableName() string {
	return "shot_frame_images"
}

// Create 创建分镜图片
func (shotFrameImages *ShotFrameImages) Create() {
	database.DB.Create(&shotFrameImages)
}

// Save 保存分镜图片
func (shotFrameImages *ShotFrameImages) Save() (rowsAffected int64) {
	result := database.DB.Save(&shotFrameImages)
	return result.RowsAffected
}

// Delete 删除分镜图片
func (shotFrameImages *ShotFrameImages) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&shotFrameImages)
	return result.RowsAffected
}
