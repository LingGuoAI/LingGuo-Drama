package shot_props

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

type ShotProps struct {
	models.BaseModel
	ShotId  uint64 `json:"shotId" form:"shotId" gorm:"column:shot_id;not null;default:0;comment:分镜ID;"`
	PropsId uint64 `json:"propsId" form:"propsId" gorm:"column:props_id;not null;default:0;comment:道具ID"`
	models.CommonTimestampsField
}

// TableName 分镜和道具中间表
func (ShotProps) TableName() string {
	return "shot_props"
}

// Create 创建分镜和道具中间
func (shotProps *ShotProps) Create() {
	database.DB.Create(&shotProps)
}

// Save 保存分镜和道具中间
func (shotProps *ShotProps) Save() (rowsAffected int64) {
	result := database.DB.Save(&shotProps)
	return result.RowsAffected
}

// Delete 删除分镜和道具中间
func (shotProps *ShotProps) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&shotProps)
	return result.RowsAffected
}
