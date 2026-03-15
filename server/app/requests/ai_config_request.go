package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// AiConfigRequest AI配置请求结构体
type AiConfigRequest struct {
	Name        string   `valid:"name" json:"name"`                 // 配置名称
	ServiceType string   `valid:"service_type" json:"service_type"` // 服务类型: text, image, video
	Provider    string   `valid:"provider" json:"provider"`         // 厂商提供商
	BaseUrl     string   `valid:"base_url" json:"base_url"`         // 接口地址
	ApiKey      string   `valid:"api_key" json:"api_key"`           // API Key
	Model       []string `json:"model"`                             // 支持的模型 (数组，govalidator 对 slice 的原生支持有限，通常依赖业务逻辑验证或 custom 验证)
	Priority    int      `valid:"priority" json:"priority"`         // 优先级
	IsActive    int8     `valid:"is_active" json:"is_active"`       // 状态 1-启用 0-禁用
}

// AiConfigSave AI配置保存时的验证规则
func AiConfigSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"name":         []string{"required", "max:255", "min:1"},
		"service_type": []string{"required", "in:text,image,video"},
		"provider":     []string{"required", "max:100"},
		"base_url":     []string{"required", "max:500"},
		"api_key":      []string{"required", "max:500"},
		//"priority":     []string{"numeric"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:配置名称为必填项",
			"min:配置名称不能为空",
			"max:配置名称长度不能超过 255 个字符",
		},
		"service_type": []string{
			"required:服务类型为必填项",
			"in:服务类型必须是 text, image 或 video",
		},
		"provider": []string{
			"required:厂商提供商为必填项",
			"max:厂商提供商长度不能超过 100 个字符",
		},
		"base_url": []string{
			"required:接口地址为必填项",
			"max:接口地址长度不能超过 500 个字符",
		},
		"api_key": []string{
			"required:API Key为必填项",
			"max:API Key长度不能超过 500 个字符",
		},
	}

	return validate(data, rules, messages)
}
