package projects

import (
	"gorm.io/gorm"
	"spiritFruit/pkg/config"
	"strings"
)

// func (projects *Projects) BeforeSave(tx *gorm.DB) (err error) {}
// func (projects *Projects) BeforeCreate(tx *gorm.DB) (err error) {}
// func (projects *Projects) AfterCreate(tx *gorm.DB) (err error) {}
// func (projects *Projects) BeforeUpdate(tx *gorm.DB) (err error) {}
// func (projects *Projects) AfterUpdate(tx *gorm.DB) (err error) {}
// func (projects *Projects) AfterSave(tx *gorm.DB) (err error) {}
// func (projects *Projects) BeforeDelete(tx *gorm.DB) (err error) {}
// func (projects *Projects) AfterDelete(tx *gorm.DB) (err error) {}
// func (projects *Projects) AfterFind(tx *gorm.DB) (err error) {}
func (projects *Projects) AfterFind(tx *gorm.DB) (err error) {
	// 处理单个图片字段: Image
	if projects.Image != nil && *projects.Image != "" {
		// 防止重复拼接
		if !strings.HasPrefix(*projects.Image, "https://") && !strings.HasPrefix(*projects.Image, "http://") {
			// 从配置获取域名
			domain := config.GetString("app.url")
			// 正确处理路径分隔符
			path := strings.TrimLeft(*projects.Image, "/")
			fullPath := domain + path
			projects.Image = &fullPath
		}
	}
	return nil
}

// BeforeUpdate 在更新前去掉URL前缀，只保存相对路径
func (projects *Projects) BeforeUpdate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: Image
	if projects.Image != nil && *projects.Image != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*projects.Image, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*projects.Image, domain)
			projects.Image = &relativePath
		} else if strings.HasPrefix(*projects.Image, "https://") || strings.HasPrefix(*projects.Image, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*projects.Image, "://"); idx != -1 {
				remaining := (*projects.Image)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					projects.Image = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}

// BeforeCreate 在创建前去掉URL前缀，只保存相对路径
func (projects *Projects) BeforeCreate(tx *gorm.DB) (err error) {
	// 处理单个图片字段: Image
	if projects.Image != nil && *projects.Image != "" {
		// 获取配置的域名
		domain := config.GetString("app.url")
		// 去掉完整URL前缀，只保留相对路径
		if strings.HasPrefix(*projects.Image, domain) {
			// 去掉域名前缀
			relativePath := strings.TrimPrefix(*projects.Image, domain)
			projects.Image = &relativePath
		} else if strings.HasPrefix(*projects.Image, "https://") || strings.HasPrefix(*projects.Image, "http://") {
			// 如果是其他域名的完整URL，提取路径部分
			if idx := strings.Index(*projects.Image, "://"); idx != -1 {
				remaining := (*projects.Image)[idx+3:]
				if pathIdx := strings.Index(remaining, "/"); pathIdx != -1 {
					relativePath := remaining[pathIdx:]
					projects.Image = &relativePath
				}
			}
		}
		// 已经是相对路径的情况，不做处理
	}
	return nil
}
