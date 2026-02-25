package characters

import (
	"spiritFruit/pkg/config"
	"strings"

	"gorm.io/gorm"
)

// func (characters *Characters) BeforeSave(tx *gorm.DB) (err error) {}
// func (characters *Characters) BeforeCreate(tx *gorm.DB) (err error) {}
// func (characters *Characters) AfterCreate(tx *gorm.DB) (err error) {}
// func (characters *Characters) BeforeUpdate(tx *gorm.DB) (err error) {}
// func (characters *Characters) AfterUpdate(tx *gorm.DB) (err error) {}
// func (characters *Characters) AfterSave(tx *gorm.DB) (err error) {}
// func (characters *Characters) BeforeDelete(tx *gorm.DB) (err error) {}
// func (characters *Characters) AfterDelete(tx *gorm.DB) (err error) {}
// func (characters *Characters) AfterFind(tx *gorm.DB) (err error) {}

func (characters *Characters) AfterFind(tx *gorm.DB) (err error) {
	// 处理单个图片字段: VisualPrompt
	if characters.AvatarUrl != nil && *characters.AvatarUrl != "" {
		// 防止重复拼接
		if !strings.HasPrefix(*characters.AvatarUrl, "https://") && !strings.HasPrefix(*characters.AvatarUrl, "http://") {
			// 从配置获取域名
			domain := config.GetString("app.url")
			// 正确处理路径分隔符
			path := strings.TrimLeft(*characters.AvatarUrl, "/")
			fullPath := domain + path
			characters.AvatarUrl = &fullPath
		}
	}
	return nil
}

// BeforeUpdate 在更新前去掉URL前缀，只保存相对路径
func (characters *Characters) BeforeUpdate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: AvataUrl
	if characters.AvatarUrl != nil && *characters.VisualPrompt != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*characters.AvatarUrl, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*characters.AvatarUrl, domain)
			characters.AvatarUrl = &relativePath
		} else if strings.HasPrefix(*characters.AvatarUrl, "https://") || strings.HasPrefix(*characters.AvatarUrl, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*characters.AvatarUrl, "://"); idx != -1 {
				remaining := (*characters.AvatarUrl)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					characters.AvatarUrl = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}

// BeforeCreate 在创建前去掉URL前缀，只保存相对路径
func (characters *Characters) BeforeCreate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: AvatarUrl
	if characters.AvatarUrl != nil && *characters.AvatarUrl != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*characters.AvatarUrl, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*characters.AvatarUrl, domain)
			characters.AvatarUrl = &relativePath
		} else if strings.HasPrefix(*characters.AvatarUrl, "https://") || strings.HasPrefix(*characters.AvatarUrl, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*characters.AvatarUrl, "://"); idx != -1 {
				remaining := (*characters.AvatarUrl)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					characters.AvatarUrl = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}
