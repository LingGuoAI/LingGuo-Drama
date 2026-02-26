package shot_frame_prompts

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

type ShotFramePrompts struct {
	models.BaseModel
	ShotId      *uint64 `json:"shotId" form:"shotId" gorm:"column:shot_id;not null;default:0;comment:分镜ID;"`
	FrameType   *string `json:"frameType" form:"frameType" gorm:"column:frame_type;comment:帧类型"` // first-首帧 last-尾帧 action-动作序列 key-关键帧
	Prompt      *string `json:"prompt" form:"prompt" gorm:"column:prompt;comment:提示词"`
	Description *string `json:"description" form:"description" gorm:"column:description;comment:描述"`
	models.CommonTimestampsField
}

// TableName 分镜帧提示词表
func (ShotFramePrompts) TableName() string {
	return "shot_frame_prompts"
}

// Create 创建分镜帧提示词
func (shotFramePrompts *ShotFramePrompts) Create() {
	database.DB.Create(&shotFramePrompts)
}

// Save 保存分镜帧提示词
func (shotFramePrompts *ShotFramePrompts) Save() (rowsAffected int64) {
	result := database.DB.Save(&shotFramePrompts)
	return result.RowsAffected
}

// Delete 删除分镜帧提示词
func (shotFramePrompts *ShotFramePrompts) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&shotFramePrompts)
	return result.RowsAffected
}
