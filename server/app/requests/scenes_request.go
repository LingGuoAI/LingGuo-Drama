package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// ScenesRequest 场景创建/更新请求
type ScenesRequest struct {
	ProjectId    uint64 `json:"projectId" valid:"projectId"` // 所属项目ID
	Name         string `json:"name" valid:"name"`           // 场景名称 (通常是 地点+时间)
	ImageUrl     string `json:"imageUrl" valid:"imageUrl"`
	Location     string `json:"location" valid:"location"`         // 地点
	Time         string `json:"time" valid:"time"`                 // 时间 (日/夜/黄昏等)
	Atmosphere   string `json:"atmosphere" valid:"atmosphere"`     // 氛围描述
	VisualPrompt string `json:"visualPrompt" valid:"visualPrompt"` // AI绘画提示词
	Status       int8   `json:"status" valid:"status"`             // 状态
}

// ScenesSave 验证场景保存/更新的规则
func ScenesSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"projectId":    []string{"required", "numeric"},
		"name":         []string{"required", "min:1", "max:100"},
		"location":     []string{"required", "max:100"},
		"time":         []string{"required", "max:50"},
		"atmosphere":   []string{"max:255"},
		"visualPrompt": []string{}, // Prompt 可以很长，甚至为空，视业务而定
		"status":       []string{"numeric", "in:1,2,3"},
	}

	messages := govalidator.MapData{
		"projectId.required": []string{"所属项目ID不能为空"},
		"projectId.numeric":  []string{"项目ID必须是数字"},
		"name.required":      []string{"场景名称不能为空"},
		"name.max":           []string{"场景名称长度不能超过100个字符"},
		"location.required":  []string{"地点不能为空"},
		"location.max":       []string{"地点长度不能超过100个字符"},
		"time.required":      []string{"时间不能为空"},
		"time.max":           []string{"时间长度不能超过50个字符"},
		"atmosphere.max":     []string{"氛围描述不能超过255个字符"},
		"status.in":          []string{"状态值无效(1-待生成 2-生成中 3-已完成)"},
	}

	return validate(data, rules, messages)
}
