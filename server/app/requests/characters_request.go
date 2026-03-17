package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// CharactersRequest 角色请求结构体
type CharactersRequest struct {
	ProjectId      uint64 `valid:"projectId" json:"projectId" `          // 所属项目ID
	Name           string `valid:"name" json:"name" `                    // 角色名
	RoleType       string `valid:"roleType" json:"roleType"`             // 角色类型: main/supporting/minor
	Gender         string `valid:"gender" json:"gender"`                 // 性别(需从appearance解析或留空)
	AgeGroup       string `valid:"ageGroup" json:"ageGroup"`             // 年龄段
	Personality    string `valid:"personality" json:"personality"`       // 性格描述
	AppearanceDesc string `valid:"appearanceDesc" json:"appearanceDesc"` // 外貌长文本描述(原appearance)
	VisualPrompt   string `valid:"visualPrompt" json:"visualPrompt"`     // AI绘画专用Prompt(从appearance提取)
	AvatarUrl      string `valid:"avatarUrl" json:"avatarUrl"`           // 头像/定妆照
	VoiceId        string `valid:"voiceId" json:"voiceId"`               // TTS音色ID
	// 关联关系字段
}

// CharactersSave 角色保存时的验证规则
func CharactersSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"projectId":      []string{"required", "numeric"},
		"name":           []string{"required", "max:100", "min:1"},
		"roleType":       []string{"max:50"},
		"gender":         []string{"max:20"},
		"ageGroup":       []string{"max:50"},
		"personality":    []string{"max:10000"},
		"appearanceDesc": []string{"max:10000"},
		"avatarUrl":      []string{"max:1024"},
		"voiceId":        []string{"max:100"},
	}

	messages := govalidator.MapData{
		"projectId": []string{
			"required:所属项目ID为必填项",
			"numeric:所属项目ID必须为数字",
		},
		"name": []string{
			"required:角色名为必填项",
			"min:角色名长度不能为空",
			"max:角色名长度不能超过 100 个字符",
		},
		"roleType": []string{
			"max:角色类型: main/supporting/minor长度不能超过 50 个字符",
		},
		"gender": []string{
			"max:性别(需从appearance解析或留空)长度不能超过 20 个字符",
		},
		"ageGroup": []string{
			"max:年龄段长度不能超过 50 个字符",
		},
		"personality": []string{
			"max:性格描述长度不能超过 10000 个字符",
		},
		"appearanceDesc": []string{
			"max:外貌长文本描述(原appearance)长度不能超过 10000 个字符",
		},
		"avatarUrl": []string{
			"max:头像/定妆照长度不能超过 1024 个字符",
		},
		"voiceId": []string{
			"max:TTS音色ID长度不能超过 100 个字符",
		},
	}

	return validate(data, rules, messages)
}
