package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

func validate(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	// 配置选项
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符
		Messages:      messages,
	}
	// 开始验证
	return govalidator.New(opts).ValidateStruct()
}

type LoginByPasswordRequest struct {
	LoginID  string `valid:"login_id" json:"login_id"`
	Password string `valid:"password" json:"password,omitempty"`
}

// LoginByPassword 验证表单，返回长度等于零即通过
func LoginByPassword(data interface{}, c *gin.Context) map[string][]string {

	rules := govalidator.MapData{
		"login_id": []string{"required", "min:3"},
		"password": []string{"required", "min:6"},
	}
	messages := govalidator.MapData{
		"login_id": []string{
			"required:登录 ID 为必填项，支持手机号、邮箱和用户名",
			"min:登录 ID 长度需大于 3",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于 6",
		},
	}

	errs := validate(data, rules, messages)

	return errs
}
