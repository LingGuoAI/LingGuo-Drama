// Package models 模型通用属性和方法
package models

import (
	"gorm.io/gorm"
	"time"

	"github.com/spf13/cast"
)

// BaseModel 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time      `gorm:"column:created_at;index;" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetStringID 获取 ID 的字符串格式
func (a BaseModel) GetStringID() string {
	return cast.ToString(a.ID)
}
