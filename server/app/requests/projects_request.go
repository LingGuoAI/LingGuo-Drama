package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// ProjectsRequest 短剧项目请求结构体
type ProjectsRequest struct {
	SerialNo      string `valid:"serialNo" json:"serialNo"`       // 业务流水号
	Title         string `valid:"title" json:"title" `            // 项目名称/短剧标题
	Description   string `valid:"description" json:"description"` // 项目简介
	Style         string `valid:"style" json:"style"`
	Status        int8   `valid:"status" json:"status" `              // 状态
	Image         string `valid:"image" json:"image"`                 // 封面图
	TotalDuration uint64 `valid:"totalDuration" json:"totalDuration"` // 总时长(秒)
	Settings      string `valid:"settings" json:"settings"`           // 生成配置快照
}

// ProjectsSave 短剧项目保存时的验证规则
func ProjectsSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"serialNo":    []string{"max:64"},
		"title":       []string{"required", "max:255", "min:1"},
		"description": []string{"max:5000"},
		//"status":        []string{"required", "numeric"},
		"image": []string{"max:1024"},
		//"totalDuration": []string{"numeric"},
		//"settings":      []string{"max:255"},
	}

	messages := govalidator.MapData{
		"serialNo": []string{
			"max:业务流水号长度不能超过 64 个字符",
		},
		"title": []string{
			"required:项目名称/短剧标题为必填项",
			"min:项目名称/短剧标题长度不能为空",
			"max:项目名称/短剧标题长度不能超过 255 个字符",
		},
		"description": []string{
			"max:项目简介长度不能超过 5000 个字符",
		},
		//"status": []string{
		//	"required:状态为必填项",
		//	"numeric:状态必须为数字",
		//},
		"image": []string{
			"max:封面图长度不能超过 1024 个字符",
		},
		//"totalDuration": []string{
		//	"numeric:总时长(秒)必须为数字",
		//},
		//"settings": []string{
		//	"max:生成配置快照长度不能超过 255 个字符",
		//},
	}

	return validate(data, rules, messages)
}
