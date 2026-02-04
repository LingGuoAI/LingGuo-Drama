package main

import (
	"context"
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"os"
	"spiritFruit/app/cmd"
	"spiritFruit/bootstrap"
	btsConfig "spiritFruit/config"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
)

func init() {
	// 加载 config 目录下的配置信息
	btsConfig.Initialize()

}

func testRedisConnection() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		Password: config.GetString("redis.password"),
		DB:       config.GetInt("redis.database_asynq"),
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	return err
}

func main() {

	// 应用的主入口，默认调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   "",
		Short: "A simple write project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {

			// 配置初始化，依赖命令行 --env 参数
			config.InitConfig(cmd.Env)

			// 初始化 Logger
			bootstrap.SetupLogger()

			// 初始化数据库
			bootstrap.SetupDB()

			// 初始化 Redis
			bootstrap.SetupRedis()

			// 测试 Redis 连接
			if err := testRedisConnection(); err != nil {
				console.Exit(fmt.Sprintf("Failed to connect to Redis with %s", err.Error()))
			}

			// 初始化缓存
			bootstrap.SetupCache()
		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
		cmd.CmdKey,
		cmd.CmdCache,
	)

	// 配置默认运行 Web 服务
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册全局参数，--env
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
