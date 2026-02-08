package props

import (
	"gorm.io/gorm"
	"spiritFruit/pkg/config"
	"strings"
)

func (props *Props) AfterFind(tx *gorm.DB) (err error) {
	// 处理单个图片字段: ImageUrl
	if props.ImageUrl != nil && *props.ImageUrl != "" {
		// 防止重复拼接
		if !strings.HasPrefix(*props.ImageUrl, "https://") && !strings.HasPrefix(*props.ImageUrl, "http://") {
			// 从配置获取域名
			domain := config.GetString("app.url")
			// 正确处理路径分隔符
			path := strings.TrimLeft(*props.ImageUrl, "/")
			fullPath := domain + path
			props.ImageUrl = &fullPath
		}
	}
	return nil
}

// BeforeUpdate 在更新前去掉URL前缀，只保存相对路径
func (props *Props) BeforeUpdate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: ImageUrl
	if props.ImageUrl != nil && *props.ImageUrl != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*props.ImageUrl, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*props.ImageUrl, domain)
			props.ImageUrl = &relativePath
		} else if strings.HasPrefix(*props.ImageUrl, "https://") || strings.HasPrefix(*props.ImageUrl, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*props.ImageUrl, "://"); idx != -1 {
				remaining := (*props.ImageUrl)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					props.ImageUrl = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}

// BeforeCreate 在创建前去掉URL前缀，只保存相对路径
func (props *Props) BeforeCreate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: ImageUrl
	if props.ImageUrl != nil && *props.ImageUrl != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*props.ImageUrl, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*props.ImageUrl, domain)
			props.ImageUrl = &relativePath
		} else if strings.HasPrefix(*props.ImageUrl, "https://") || strings.HasPrefix(*props.ImageUrl, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*props.ImageUrl, "://"); idx != -1 {
				remaining := (*props.ImageUrl)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					props.ImageUrl = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}
