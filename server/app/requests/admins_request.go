package requests

import (
    "github.com/gin-gonic/gin"
    "github.com/thedevsaddam/govalidator"
)

// AdminsRequest 系统管理员请求结构体
type AdminsRequest struct {
    Username string `valid:"username" json:"username" ` // 用户名
    Mobile string `valid:"mobile" json:"mobile" ` // 手机号
    Password string `valid:"password" json:"password" ` // 密码
    Email string `valid:"email" json:"email"` // 邮箱
    AuthorityId uint64 `valid:"authorityId" json:"authorityId"` // 用户角色ID
}

// AdminsSave 系统管理员保存时的验证规则
func AdminsSave(data interface{}, c *gin.Context) map[string][]string {
    rules := govalidator.MapData{
        "username": []string{"required", "max:120", "min:1"},
        "mobile": []string{"required", "max:11", "min:1"},
        "password": []string{"required", "max:64", "min:1"},
        "email": []string{"max:80"},
        "authorityId": []string{"numeric"},
    }

    messages := govalidator.MapData{
        "username": []string{
             "required:用户名为必填项",
             "min:用户名长度不能为空",
             "max:用户名长度不能超过 120 个字符",
        },
        "mobile": []string{
             "required:手机号为必填项",
             "min:手机号长度不能为空",
             "max:手机号长度不能超过 11 个字符",
        },
        "password": []string{
             "required:密码为必填项",
             "min:密码长度不能为空",
             "max:密码长度不能超过 64 个字符",
        },
        "email": []string{
             "max:邮箱长度不能超过 80 个字符",
        },
        "authorityId": []string{
             "numeric:用户角色ID必须为数字",
        },
    }

    return validate(data, rules, messages)
}