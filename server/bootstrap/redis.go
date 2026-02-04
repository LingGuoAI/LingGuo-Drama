package bootstrap

import (
	"fmt"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/redis"
)

// SetupRedis 初始化 Redis
func SetupRedis() {

	// 建立 Redis 连接
	redis.ConnectRedis(
		fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		config.GetString("redis.username"),
		config.GetString("redis.password"),
		config.GetInt("redis.database"),
	)
}
