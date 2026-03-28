package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type RegisterRequest struct {
	Username string `valid:"username" json:"username"`
	Mobile   string `valid:"mobile" json:"mobile"`
	Email    string `valid:"email" json:"email"`
	Password string `valid:"password" json:"password"`
}

// Register 验证注册表单
func Register(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"username": []string{"required", "alpha_num", "between:3,20", "not_exists:admins,username"},
		"mobile":   []string{"required", "digits:11", "not_exists:admins,mobile"},
		"email":    []string{"required", "email", "not_exists:admins,email"},
		"password": []string{"required", "min:6"},
	}
	messages := govalidator.MapData{
		"username": []string{
			"required:用户名必填",
			"alpha_num:用户名格式不正确，仅支持字母和数字",
			"between:用户名长度需在 3-20 之间",
			"not_exists:用户名已被占用",
		},
		"mobile": []string{
			"required:手机号必填",
			"digits:手机号格式不正确，需为 11 位数字",
			"not_exists:手机号已被占用",
		},
		"email": []string{
			"required:邮箱必填",
			"email:邮箱格式不正确",
			"not_exists:邮箱已被占用",
		},
		"password": []string{
			"required:密码必填",
			"min:密码长度需大于 6",
		},
	}

	errs := validate(data, rules, messages)

	return errs
}
