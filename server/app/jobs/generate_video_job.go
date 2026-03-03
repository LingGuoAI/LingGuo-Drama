// app/jobs/generate_video_job.go
package jobs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/shots"
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/upload"
	"spiritFruit/pkg/video" // 引入您前面封装的 video 包
)

// VideoAPIConfig 准备视频 AI 配置结构体
type VideoAPIConfig struct {
	GetGoBaseURL   string
	GetGoKey       string
	VolcesBaseURL  string
	VolcesKey      string
	MinimaxBaseURL string
	MinimaxKey     string
	RunwayBaseURL  string
	RunwayKey      string
	PikaBaseURL    string
	PikaKey        string
	OpenAIBaseURL  string
	OpenAIKey      string
}

// helper: 根据模型名字自动获取对应的 API 配置
func getProviderConfig(modelName string) (provider, baseURL, apiKey string) {
	// 统一获取配置数据
	cfg := VideoAPIConfig{
		GetGoBaseURL:   config.GetString("ai.getgoapi.base_url"),
		GetGoKey:       config.GetString("ai.getgoapi.api_key"),
		VolcesBaseURL:  config.GetString("ai.volces.base_url"),
		VolcesKey:      config.GetString("ai.volces.api_key"),
		MinimaxBaseURL: config.GetString("ai.minimax.base_url"),
		MinimaxKey:     config.GetString("ai.minimax.api_key"),
		RunwayBaseURL:  config.GetString("ai.runway.base_url"),
		RunwayKey:      config.GetString("ai.runway.api_key"),
		PikaBaseURL:    config.GetString("ai.pika.base_url"),
		PikaKey:        config.GetString("ai.pika.api_key"),
		OpenAIBaseURL:  config.GetString("ai.openai.base_url"),
		OpenAIKey:      config.GetString("ai.openai.api_key"),
	}

	modelName = strings.ToLower(modelName)

	if strings.Contains(modelName, "doubao") || strings.Contains(modelName, "seedance") {
		return "volces", cfg.VolcesBaseURL, cfg.VolcesKey
	} else if strings.Contains(modelName, "sora") {
		return "openai", cfg.OpenAIBaseURL, cfg.OpenAIKey
	} else if strings.Contains(modelName, "runway") {
		return "runway", cfg.RunwayBaseURL, cfg.RunwayKey
	} else if strings.Contains(modelName, "pika") {
		return "pika", cfg.PikaBaseURL, cfg.PikaKey
	} else if strings.Contains(modelName, "minimax") || strings.Contains(modelName, "hailuo") {
		return "minimax", cfg.MinimaxBaseURL, cfg.MinimaxKey
	}

	// 默认使用中转平台兜底 (比如 kling 等未特定归类的模型)
	return "getgoapi", cfg.GetGoBaseURL, cfg.GetGoKey
}

// HandleGenerateVideoTask 处理视频生成任务 (包含调用与轮询)
func HandleGenerateVideoTask(ctx context.Context, t *asynq.Task) error {
	var p myAsynq.GenerateVideoPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 1. 获取并标记任务开始
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		return nil
	}
	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始生成视频 (ShotID: %d, Model: %s)", p.AsyncTaskID, p.ShotID, p.Model))

	// 2. 初始化客户端
	taskModel.UpdateProgress(10)
	provider, baseURL, apiKey := getProviderConfig(p.Model)
	// 实例化 Client
	client, err := video.NewClient(provider, baseURL, apiKey, p.Model, "", "")
	if err != nil {
		taskModel.MarkAsFailed(err)
		return err
	}

	// 3. 构建参数选项
	var opts []video.VideoOption
	opts = append(opts, video.WithDuration(p.Duration))

	appURL := config.GetString("app.url")

	fixURL := func(url string) string {
		if url == "" || strings.HasPrefix(url, "http") || strings.HasPrefix(url, "data:") {
			return url
		}

		cleanPath := strings.TrimPrefix(url, "/")

		fileData, err := os.ReadFile(cleanPath)
		if err == nil {
			// 读取成功，转换成 base64 Data URI
			ext := filepath.Ext(cleanPath)
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "image/jpeg" // 默认 mime
			}
			base64Str := base64.StdEncoding.EncodeToString(fileData)
			return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
		}

		console.Warning("读取本地图片转Base64失败，回退到URL拼接: " + err.Error())
		return strings.TrimRight(appURL, "/") + "/" + cleanPath
	}

	if p.ImageURL != "" {
		opts = append(opts, video.WithImageURL(fixURL(p.ImageURL)))
	}
	if p.FirstFrameURL != "" {
		opts = append(opts, video.WithFirstFrame(fixURL(p.FirstFrameURL)))
	}
	if p.LastFrameURL != "" {
		opts = append(opts, video.WithLastFrame(fixURL(p.LastFrameURL)))
	}
	if len(p.ReferenceImageURLs) > 0 {
		var fixedURLs []string
		for _, u := range p.ReferenceImageURLs {
			fixedURLs = append(fixedURLs, fixURL(u))
		}
		opts = append(opts, video.WithReferenceImages(fixedURLs))
	}

	// 4. 发起生成请求
	taskModel.UpdateProgress(30)
	result, err := client.GenerateVideo(p.Prompt, opts...)
	if err != nil {
		taskModel.MarkAsFailed(err)
		return err
	}

	console.Success(fmt.Sprintf("任务[%d] - 视频请求已提交，TaskID: %s, 进入轮询...", p.AsyncTaskID, result.TaskID))

	// 5. 轮询获取任务结果 (视频生成较慢，通常需 1~5 分钟)
	taskModel.UpdateProgress(40)
	maxAttempts := 150 // 最大重试 150 次，每次 10 秒 = 25分钟超时
	interval := 10 * time.Second

	for attempt := 0; attempt < maxAttempts; attempt++ {
		time.Sleep(interval)

		// 检查数据库任务是否被人工终止
		var checkTask async_tasks.AsyncTask
		if err := database.DB.First(&checkTask, p.AsyncTaskID).Error; err == nil {
			if checkTask.Status != async_tasks.StatusProcessing {
				return nil // 任务状态已变，终止轮询
			}
		}

		// 轮询远端 API
		statusRes, err := client.GetTaskStatus(result.TaskID)
		if err != nil {
			console.Warning(fmt.Sprintf("任务[%d] 轮询失败: %v，继续重试...", p.AsyncTaskID, err))
			continue
		}

		if statusRes.Error != "" {
			errStr := fmt.Errorf("video generation failed: %s", statusRes.Error)
			taskModel.MarkAsFailed(errStr)
			return errStr
		}

		if statusRes.Completed && statusRes.VideoURL != "" {
			// 视频生成成功，退出轮询
			result = statusRes
			break
		}

		// 更新进度
		prog := 40 + int(float64(attempt)/float64(maxAttempts)*45)
		taskModel.UpdateProgress(uint64(prog))
	}

	if !result.Completed || result.VideoURL == "" {
		errStr := fmt.Errorf("视频生成超时或未返回 URL")
		taskModel.MarkAsFailed(errStr)
		return errStr
	}

	console.Success(fmt.Sprintf("任务[%d] - 视频生成完成，开始下载...", p.AsyncTaskID))

	// 6. 下载视频并保存到本地
	taskModel.UpdateProgress(90)
	localPath, err := upload.DownloadAndSaveVideo(result.VideoURL)
	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("save video failed: %v", err))
		return err
	}

	// 7. 更新 Shots 镜头表中的视频链接和生成时长
	err = database.DB.Model(&shots.Shots{}).
		Where("id = ?", p.ShotID).
		Updates(map[string]interface{}{
			"video_url":   localPath,
			"duration_ms": p.Duration * 1000,
		}).Error

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("db update failed: %v", err))
		return err
	}

	// 8. 标记任务成功
	taskModel.UpdateProgress(100)
	resData := map[string]interface{}{
		"url":      localPath,
		"shot_id":  p.ShotID,
		"duration": p.Duration,
	}
	resBytes, _ := json.Marshal(resData)
	taskModel.MarkAsSuccess(string(resBytes))

	console.Success(fmt.Sprintf("任务[%d] - 视频任务彻底完成", p.AsyncTaskID))
	return nil
}
