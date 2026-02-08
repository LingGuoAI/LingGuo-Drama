package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type PropsRequest struct {
	ProjectId   uint64 `json:"projectId" valid:"projectId"`
	Name        string `json:"name" valid:"name"`
	Type        string `json:"type" valid:"type"`
	Description string `json:"description" valid:"description"`
	ImagePrompt string `json:"imagePrompt" valid:"imagePrompt"`
	ImageUrl    string `json:"imageUrl" valid:"imageUrl"`
}

func PropsSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"projectId":   []string{"required", "numeric"},
		"name":        []string{"required", "min:1", "max:100"},
		"type":        []string{"max:50"},
		"description": []string{},
		"imagePrompt": []string{},
		"imageUrl":    []string{},
	}
	messages := govalidator.MapData{
		"projectId": []string{
			"required:项目ID为必填项",
			"numeric:项目ID必须是数字",
		},
		"name": []string{
			"required:道具名称为必填项",
			"min:道具名称长度需至少 1 个字符",
			"max:道具名称长度不能超过 100 个字符",
		},
	}
	return validate(data, rules, messages)
}
