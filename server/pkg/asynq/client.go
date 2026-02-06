package asynq

import (
	"fmt"
	"github.com/hibiken/asynq"
	"spiritFruit/pkg/config"
	"sync"
)

var (
	client *asynq.Client
	once   sync.Once
)

// GetClient GetClient初始化并返回ASYNQ的Singleton实例。客户端
func GetClient() *asynq.Client {
	once.Do(func() {
		client = asynq.NewClient(asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
			Password: config.GetString("redis.password"),
			DB:       config.GetInt("redis.database_async"),
		})
	})
	return client
}
