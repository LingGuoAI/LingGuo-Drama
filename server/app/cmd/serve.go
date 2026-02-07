package cmd

import (
	"context"
	"net/http"
	"spiritFruit/bootstrap"
	"spiritFruit/pkg/appctx"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs,
}

func runWeb(cmd *cobra.Command, args []string) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	bootstrap.SetupRoute(router)

	// 使用 http.Server 替代 router.Run()
	srv := &http.Server{
		Addr:    ":" + config.Get("app.port"),
		Handler: router,
	}

	// 在 goroutine 中启动服务器
	go func() {
		console.Success("Server started on port " + config.Get("app.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorString("CMD", "serve", err.Error())
			console.Exit("Unable to start server, error:" + err.Error())
		}
	}()

	// 等待全局 context 被取消
	<-appctx.GetContext().Done()

	console.Warning("Shutting down server...")

	// 创建一个带超时的 context 用于关闭（给正在处理的请求一些时间）
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.ErrorString("CMD", "serve", "Server forced to shutdown: "+err.Error())
	}

	console.Success("Server exited properly")
}
