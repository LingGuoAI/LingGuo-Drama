package source

import (
	"gorm.io/gorm"
	"spiritFruit/pkg/config"
	"strings"
)

// func (source *Source) BeforeSave(tx *gorm.DB) (err error) {}
// func (source *Source) BeforeCreate(tx *gorm.DB) (err error) {}
// func (source *Source) AfterCreate(tx *gorm.DB) (err error) {}
// func (source *Source) BeforeUpdate(tx *gorm.DB) (err error) {}
// func (source *Source) AfterUpdate(tx *gorm.DB) (err error) {}
// func (source *Source) AfterSave(tx *gorm.DB) (err error) {}
// func (source *Source) BeforeDelete(tx *gorm.DB) (err error) {}
// func (source *Source) AfterDelete(tx *gorm.DB) (err error) {}
// func (source *Source) AfterFind(tx *gorm.DB) (err error) {}

func (source *Source) AfterFind(tx *gorm.DB) (err error) {
	// 处理单个图片字段: Image
	if source.VideoUrl != nil && *source.VideoUrl != "" {
		// 防止重复拼接
		if !strings.HasPrefix(*source.VideoUrl, "https://") && !strings.HasPrefix(*source.VideoUrl, "http://") {
			// 从配置获取域名
			domain := config.GetString("app.url")
			// 正确处理路径分隔符
			path := strings.TrimLeft(*source.VideoUrl, "/")
			fullPath := domain + path
			source.VideoUrl = &fullPath
		}
	}
	return nil
}

// BeforeUpdate 在更新前去掉URL前缀，只保存相对路径
func (source *Source) BeforeUpdate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: Image
	if source.VideoUrl != nil && *source.VideoUrl != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*source.VideoUrl, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*source.VideoUrl, domain)
			source.VideoUrl = &relativePath
		} else if strings.HasPrefix(*source.VideoUrl, "https://") || strings.HasPrefix(*source.VideoUrl, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*source.VideoUrl, "://"); idx != -1 {
				remaining := (*source.VideoUrl)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					source.VideoUrl = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}

// BeforeCreate 在创建前去掉URL前缀，只保存相对路径
func (source *Source) BeforeCreate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: Image
	if source.VideoUrl != nil && *source.VideoUrl != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*source.VideoUrl, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*source.VideoUrl, domain)
			source.VideoUrl = &relativePath
		} else if strings.HasPrefix(*source.VideoUrl, "https://") || strings.HasPrefix(*source.VideoUrl, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*source.VideoUrl, "://"); idx != -1 {
				remaining := (*source.VideoUrl)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					source.VideoUrl = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}
