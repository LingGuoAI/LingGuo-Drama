package main

import (
	"context"
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"spiritFruit/app/cmd"
	"spiritFruit/bootstrap"
	btsConfig "spiritFruit/config"
	"spiritFruit/pkg/appctx"
	"spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"syscall"
)

func init() {
	// 加载 config 目录下的配置信息
	btsConfig.Initialize()

}

// setupGracefulShutdown 设置优雅关闭
func setupGracefulShutdown() {
	// 初始化全局 context
	ctx, cancel := appctx.Initialize()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigs
		console.Warning(fmt.Sprintf("Received signal: %v, initiating graceful shutdown...", sig))

		cancel()
	}()

	_ = ctx // 避免未使用警告
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
	setupGracefulShutdown()
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

			asynq.GetClient()

			// 使用共享的 context
			go asynq.SetAsyncServerClient(appctx.GetContext())
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
