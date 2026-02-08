package scenes

import (
	"gorm.io/gorm"
	"spiritFruit/pkg/config"
	"strings"
)

func (scenes *Scenes) AfterFind(tx *gorm.DB) (err error) {
	// 处理单个图片字段: VisualPrompt
	if scenes.VisualPrompt != nil && *scenes.VisualPrompt != "" {
		// 防止重复拼接
		if !strings.HasPrefix(*scenes.VisualPrompt, "https://") && !strings.HasPrefix(*scenes.VisualPrompt, "http://") {
			// 从配置获取域名
			domain := config.GetString("app.url")
			// 正确处理路径分隔符
			path := strings.TrimLeft(*scenes.VisualPrompt, "/")
			fullPath := domain + path
			scenes.VisualPrompt = &fullPath
		}
	}
	return nil
}

// BeforeUpdate 在更新前去掉URL前缀，只保存相对路径
func (scenes *Scenes) BeforeUpdate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: VisualPrompt
	if scenes.VisualPrompt != nil && *scenes.VisualPrompt != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*scenes.VisualPrompt, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*scenes.VisualPrompt, domain)
			scenes.VisualPrompt = &relativePath
		} else if strings.HasPrefix(*scenes.VisualPrompt, "https://") || strings.HasPrefix(*scenes.VisualPrompt, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*scenes.VisualPrompt, "://"); idx != -1 {
				remaining := (*scenes.VisualPrompt)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					scenes.VisualPrompt = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}

// BeforeCreate 在创建前去掉URL前缀，只保存相对路径
func (scenes *Scenes) BeforeCreate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: VisualPrompt
	if scenes.VisualPrompt != nil && *scenes.VisualPrompt != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*scenes.VisualPrompt, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*scenes.VisualPrompt, domain)
			scenes.VisualPrompt = &relativePath
		} else if strings.HasPrefix(*scenes.VisualPrompt, "https://") || strings.HasPrefix(*scenes.VisualPrompt, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*scenes.VisualPrompt, "://"); idx != -1 {
				remaining := (*scenes.VisualPrompt)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					scenes.VisualPrompt = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}
