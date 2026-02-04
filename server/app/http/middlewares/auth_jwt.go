// Package middlewares Gin 中间件
package middlewares

import (
	"fmt"
	"spiritFruit/app/models/admins"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/jwt"
	"spiritFruit/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthAdminJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 从标头 Authorization:Bearer xxxxx 中获取信息，并验证 JWT 的准确性
		claims, err := jwt.NewJWT().ParserToken(c)

		// JWT 解析失败，有错误发生
		if err != nil {
			response.Unauthorized(c, fmt.Sprintf("请查看 %v 相关的接口认证文档", config.GetString("app.name")))
			return
		}

		// JWT 解析成功，设置用户信息
		adminModel := admins.Get(claims.UserID)
		if adminModel.ID == 0 {
			response.Unauthorized(c, "找不到对应用户，用户可能已删除")
			return
		}

		// 将用户信息存入 gin.context 里，后续 auth 包将从这里拿到当前用户数据
		c.Set("current_admin_id", adminModel.GetStringID())
		c.Set("current_admin_name", adminModel.Username)
		c.Set("current_admin", adminModel)

		c.Next()
	}
}
