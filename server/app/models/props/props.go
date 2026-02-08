package props

import (
	"spiritFruit/app/models"
	"spiritFruit/app/models/projects"
	"spiritFruit/pkg/database"
)

// Props 道具表
type Props struct {
	models.BaseModel

	// 基础信息
	ProjectId   *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;comment:所属项目ID;index"`
	Name        *string `json:"name" form:"name" gorm:"column:name;comment:道具名称;size:100;"`
	Type        *string `json:"type" form:"type" gorm:"column:type;comment:道具类型(如:交通工具/武器/电子产品);size:50;"`
	Description *string `json:"description" form:"description" gorm:"column:description;comment:道具描述(剧情描述);type:text;"`

	// AI 生图相关
	ImagePrompt *string `json:"imagePrompt" form:"imagePrompt" gorm:"column:image_prompt;comment:AI绘画提示词(Visual Prompt);type:text;"`
	ImageUrl    *string `json:"imageUrl" form:"imageUrl" gorm:"column:image_url;comment:道具图片URL;size:1024;"`

	// 关联关系
	Project *projects.Projects `json:"project,omitempty" gorm:"foreignKey:ProjectId;references:ID"`

	models.CommonTimestampsField
}

// TableName 指定表名
func (Props) TableName() string {
	return "props"
}

// Create 创建道具项目
func (props *Props) Create() {
	database.DB.Create(&props)
}

// Save 保存道具项目
func (props *Props) Save() (rowsAffected int64) {
	result := database.DB.Save(&props)
	return result.RowsAffected
}

// Delete 删除道具项目
func (props *Props) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&props)
	return result.RowsAffected
}
