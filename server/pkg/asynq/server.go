package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
)

var server *asynq.Server

// SetAsyncServerClient 以并发和队列选项初始化服务器
func SetAsyncServerClient(ctx context.Context) {
	redisConnOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		Password: config.GetString("redis.password"),
		DB:       config.GetInt("redis.database_asynq"),
	}

	mux := asynq.NewServeMux()
	mux.HandleFunc("sora:generate_video", handelGenerateVideo)

	server = asynq.NewServer(
		redisConnOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default":  6,
				"critical": 4,
			},
		},
	)

	// 监听 context 取消，优雅关闭
	go func() {
		<-ctx.Done()
		console.Success("Shutting down asynq server...")
		server.Shutdown()
	}()

	if err := server.Run(mux); err != nil {
		// 如果是因为 shutdown 导致的退出，不算错误
		if ctx.Err() != nil {
			console.Success("Asynq server stopped gracefully")
			return
		}
		console.Exit(fmt.Sprintf("Failed to run asynq server with %s", err.Error()))
	}
}

// Shutdown 提供外部调用的关闭方法
func Shutdown() {
	if server != nil {
		server.Shutdown()
	}
}

type GenerateVideoPaperPayload struct {
	OrderSn string `json:"order_sn"`
}

func handelGenerateVideo(ctx context.Context, t *asynq.Task) error {
	var p GenerateVideoPaperPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		console.Error(fmt.Sprintf("Failed to unmarshal payload: %v", err))
		return err
	}

	return fmt.Errorf("video generation timeout")
}
