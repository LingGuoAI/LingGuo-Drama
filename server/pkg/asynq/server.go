package asynq

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
)

var server *asynq.Server

// StartServer 启动服务
func StartServer(ctx context.Context, mux *asynq.ServeMux) {
	// 1. 初始化 Redis 连接配置
	redisConnOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		Password: config.GetString("redis.password"),
		DB:       config.GetInt("redis.database_async"),
	}

	// 2. 创建 Server 实例
	server = asynq.NewServer(
		redisConnOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				console.Error(fmt.Sprintf("任务处理失败 [Type: %s]: %v", task.Type(), err))
			}),
		},
	)

	// 3. 监听 Context 取消信号
	go func() {
		<-ctx.Done()
		console.Success("正在关闭 Asynq Server...")
		Shutdown()
	}()

	// 4. 启动服务
	console.Success("Asynq Server 已启动，正在监听任务...")
	if err := server.Run(mux); err != nil {
		if ctx.Err() != nil {
			console.Success("Asynq Server 已优雅停止")
			return
		}
		console.Exit(fmt.Sprintf("Asynq Server 启动失败: %s", err.Error()))
	}
}

// Shutdown 提供外部调用的关闭方法
func Shutdown() {
	if server != nil {
		server.Shutdown()
	}
}
