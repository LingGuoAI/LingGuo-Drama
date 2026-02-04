// Package auth 授权相关逻辑
package auth

import (
	"errors"
	"spiritFruit/app/models/admins"
	"spiritFruit/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AdminAttempt 代理尝试登录
func AdminAttempt(value string, password string) (admins.Admins, error) {
	adminModel := admins.GetByMulti(value)
	if adminModel.ID == 0 {
		return admins.Admins{}, errors.New("账号不存在")
	}

	if !adminModel.ComparePassword(password) {
		return admins.Admins{}, errors.New("密码错误")
	}

	return adminModel, nil
}

// CurrentAdmin 获取当前登录代理信息
func CurrentAdmin(c *gin.Context) admins.Admins {
	adminsModel, ok := c.MustGet("current_admin").(admins.Admins)
	if !ok {
		logger.LogIf(errors.New("无法获取代理信息"))
		return admins.Admins{}
	}
	return adminsModel
}

// CurrentUID 从 gin.context 中获取当前登录用户 ID
func CurrentUID(c *gin.Context) string {
	return c.GetString("current_user_id")
}
