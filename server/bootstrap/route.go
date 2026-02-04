// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
	"net/http"
	"spiritFruit/app/http/middlewares"
	"spiritFruit/routes"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoute 路由初始化
func SetupRoute(router *gin.Engine) {

	// 注册全局中间件
	registerGlobalMiddleWare(router)

	// 注册 Swagger 处理器
	registerSwagger(router)

	// 注册管理端路由
	routes.RegisterAdminAPIRoutes(router)

	//  配置 404 路由
	setup404Handler(router)
}

func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(
		middlewares.Logger(),
		middlewares.Recovery(),
		middlewares.ForceUA(),
		middlewares.Cors(),
	)
}

// 注册 Swagger 处理器
func registerSwagger(router *gin.Engine) {
	// 访问路径: http://localhost:8080/swagger/index.html
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func setup404Handler(router *gin.Engine) {
	// 处理 404 请求
	router.NoRoute(func(c *gin.Context) {
		// 获取标头信息的 Accept 信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 如果是 HTML 的话
			c.String(http.StatusNotFound, "页面返回 404")
		} else {
			// 默认返回 JSON
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
			})
		}
	})
}
