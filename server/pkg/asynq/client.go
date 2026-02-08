package asynq

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"sync"
)

var (
	client *asynq.Client
	once   sync.Once
)

// GetClient 获取/初始化 Asynq Client 单例
func GetClient() *asynq.Client {
	once.Do(func() {
		redisOpt := asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
			Password: config.GetString("redis.password"),
			DB:       config.GetInt("redis.database_asynq"),
		}
		client = asynq.NewClient(redisOpt)
	})
	return client
}

// EnqueueGenerateScript 投递剧本生成任务
func EnqueueGenerateScript(payload GenerateScriptPayload) (*asynq.TaskInfo, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(TypeGenerateScript, bytes)
	return GetClient().Enqueue(task, asynq.Queue("critical"))
}

// EnqueueGenerateImage 投递图片生成任务
func EnqueueGenerateImage(payload GenerateImagePayload) (*asynq.TaskInfo, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(TypeGenerateImage, bytes)
	// 图片生成可能较慢，放入 default 队列
	info, err := GetClient().Enqueue(task, asynq.Queue("default"))
	if err != nil {
		console.Error(fmt.Sprintf("投递图片任务失败: %v", err))
		return nil, err
	}
	return info, nil
}

// EnqueueGenerateCharacters 投递角色生成任务
func EnqueueGenerateCharacters(payload GenerateCharactersPayload) (*asynq.TaskInfo, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// 角色生成不算特别耗时，放入 default 队列
	task := asynq.NewTask(TypeGenerateCharacters, bytes)
	return GetClient().Enqueue(task, asynq.Queue("default"))
}

// EnqueueExtractScenes 投递场景提取任务
func EnqueueExtractScenes(payload ExtractScenesPayload) (*asynq.TaskInfo, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// 场景提取涉及整集剧本分析，耗时较长，放入 default 队列
	task := asynq.NewTask(TypeExtractScenes, bytes)
	return GetClient().Enqueue(task, asynq.Queue("default"))
}

// EnqueueGenerateSceneImage 投递场景生图任务
func EnqueueGenerateSceneImage(payload GenerateSceneImagePayload) (*asynq.TaskInfo, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// 放入 default 队列
	task := asynq.NewTask(TypeGenerateSceneImage, bytes)
	return GetClient().Enqueue(task, asynq.Queue("default"))
}

// EnqueueGenerateShots 投递分镜生成任务
func EnqueueGenerateShots(payload GenerateShotsPayload) (*asynq.TaskInfo, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	// 分镜生成耗时较长，建议放入 default 或 critical 队列
	task := asynq.NewTask(TypeGenerateShots, bytes)
	return GetClient().Enqueue(task, asynq.Queue("default"))
}
