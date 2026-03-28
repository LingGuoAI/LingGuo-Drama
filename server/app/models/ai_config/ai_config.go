package ai_config

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

// AiConfig 结构体 AI服务配置
type AiConfig struct {
	models.BaseModel
	Name        *string  `json:"name" form:"name" gorm:"column:name;comment:配置名称;size:255;"`
	ServiceType *string  `json:"service_type" form:"service_type" gorm:"column:service_type;comment:服务类型(text, image, video);size:50;"`
	Provider    *string  `json:"provider" form:"provider" gorm:"column:provider;comment:服务提供商;size:100;"`
	BaseUrl     *string  `json:"base_url" form:"base_url" gorm:"column:base_url;comment:接口地址;size:500;"`
	ApiKey      *string  `json:"api_key" form:"api_key" gorm:"column:api_key;comment:API Key;size:500;"`
	Model       []string `json:"model" form:"model" gorm:"column:model;type:json;serializer:json;comment:支持的模型列表;"` // 自动序列化为JSON
	Priority    *int     `json:"priority" form:"priority" gorm:"default:0;column:priority;comment:优先级;"`
	IsActive    *int8    `json:"is_active" form:"is_active" gorm:"default:1;column:is_active;comment:状态 1-启用 0-禁用;"`
	AdminID     *uint64  `json:"admin_id" form:"admin_id" gorm:"column:admin_id;index;comment:所属管理员ID;"`
	models.CommonTimestampsField
}

// TableName 自定义表名 ai_config
func (AiConfig) TableName() string {
	return "ai_config"
}

// Create 创建AI配置
func (config *AiConfig) Create() {
	database.DB.Create(&config)
}

// Save 保存AI配置
func (config *AiConfig) Save() (rowsAffected int64) {
	result := database.DB.Save(&config)
	return result.RowsAffected
}

// Delete 删除AI配置
func (config *AiConfig) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&config)
	return result.RowsAffected
}
