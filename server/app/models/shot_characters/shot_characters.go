package shot_characters

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

type ShotCharacters struct {
	models.BaseModel
	ShotId      uint64 `json:"shotId" form:"shotId" gorm:"column:shot_id;not null;default:0;comment:分镜ID;"`
	CharacterId uint64 `json:"characterId" form:"characterId" gorm:"column:character_id;not null; default:0;comment:角色ID"`
	models.CommonTimestampsField
}

// TableName 分镜和角色中间表
func (ShotCharacters) TableName() string {
	return "shot_characters"
}

// Create 创建分镜和角色中间
func (shotCharacters *ShotCharacters) Create() {
	database.DB.Create(&shotCharacters)
}

// Save 保存分镜和角色中间
func (shotCharacters *ShotCharacters) Save() (rowsAffected int64) {
	result := database.DB.Save(&shotCharacters)
	return result.RowsAffected
}

// Delete 删除分镜和角色中间
func (shotCharacters *ShotCharacters) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&shotCharacters)
	return result.RowsAffected
}
