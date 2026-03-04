package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"

	"spiritFruit/app/cmd"
	"spiritFruit/app/jobs"
	"spiritFruit/bootstrap"
	btsConfig "spiritFruit/config"
	"spiritFruit/pkg/appctx"
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
)

func init() {
	// 加载 config 目录下的配置映射
	btsConfig.Initialize()
}

func main() {

	// 设置优雅关闭信号监听
	setupGracefulShutdown()
	// 初始化应用命令
	var rootCmd = &cobra.Command{
		Use:   "spiritFruit",
		Short: "An AI-powered backend for short drama script and image generation",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		// PersistentPreRun 在任何子命令执行前都会运行
		PersistentPreRun: func(command *cobra.Command, args []string) {

			// 1. 初始化配置和基础设施
			config.InitConfig(cmd.Env)
			bootstrap.SetupLogger()
			bootstrap.SetupDB()
			bootstrap.SetupAutoMigrate()
			bootstrap.SetupRedis()
			bootstrap.SetupCache()

			// 2. 初始化 Asynq 客户端 (用于投递任务)
			myAsynq.GetClient()

			mux := asynq.NewServeMux()
			mux.HandleFunc(myAsynq.TypeGenerateScript, jobs.HandleGenerateScript)
			mux.HandleFunc(myAsynq.TypeGenerateImage, jobs.HandleGenerateImage)
			mux.HandleFunc(myAsynq.TypeGenerateCharacters, jobs.HandleGenerateCharacters)
			mux.HandleFunc(myAsynq.TypeExtractScenes, jobs.HandleExtractScenes)
			mux.HandleFunc(myAsynq.TypeGenerateSceneImage, jobs.HandleGenerateSceneImage)
			mux.HandleFunc(myAsynq.TypeGenerateShots, jobs.HandleGenerateShots)
			mux.HandleFunc(myAsynq.TypeExtractProps, jobs.HandleExtractPropsTask)
			mux.HandleFunc(myAsynq.TypeGeneratePropImage, jobs.HandleGeneratePropImageTask)
			mux.HandleFunc(myAsynq.TypeExtractFramePrompt, jobs.HandleExtractFramePromptTask)
			mux.HandleFunc(myAsynq.TypeGenerateFrameImage, jobs.HandleGenerateFrameImageTask)
			mux.HandleFunc(myAsynq.TypeGenerateVideo, jobs.HandleGenerateVideoTask)
			mux.HandleFunc(myAsynq.TypeMergeVideo, jobs.HandleMergeVideoTask)

			// 4. 启动 Asynq Server (消费者)
			// 使用 goroutine 启动，避免阻塞主线程（Web Server）
			// 传入 appctx 获取的全局 context，用于优雅关闭
			go myAsynq.StartServer(appctx.GetContext(), mux)
		},

		// PersistentPostRun 在命令执行结束后运行
		PersistentPostRun: func(command *cobra.Command, args []string) {
			// 优雅关闭 Asynq Server
			myAsynq.Shutdown()
			console.Success("Application shutdown complete")
		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
		cmd.CmdKey,
		cmd.CmdCache,
	)

	// 注册默认命令 (serve)
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册全局参数 (--env)
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}

// setupGracefulShutdown 监听系统信号，触发 context 取消
func setupGracefulShutdown() {
	// 这里的 Initialize 返回的是全局 context 的 cancel 函数
	_, cancel := appctx.Initialize()

	go func() {
		// 创建信号通道
		sigs := make(chan os.Signal, 1)
		// 监听 SIGINT (Ctrl+C) 和 SIGTERM (kill)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		// 阻塞直到收到信号
		sig := <-sigs
		console.Warning(fmt.Sprintf("Received signal: %v, initiating graceful shutdown...", sig))

		// 调用 cancel，通知所有通过 appctx.GetContext() 获取 context 的组件（如 Asynq Server）停止工作
		cancel()
	}()
}
