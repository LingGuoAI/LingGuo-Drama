package video

import (
	"fmt"
	"net/http"
	"time"
)

// VideoClient 定义视频生成客户端接口
type VideoClient interface {
	// GenerateVideo 发起生成视频的请求
	GenerateVideo(prompt string, opts ...VideoOption) (*VideoResult, error)
	// GetTaskStatus 查询视频生成任务的状态
	GetTaskStatus(taskID string) (*VideoResult, error)
}

// VideoResult 表示视频生成的结果或状态
type VideoResult struct {
	TaskID       string
	Status       string // pending, processing, completed, failed
	VideoURL     string
	ThumbnailURL string
	Duration     int
	Width        int
	Height       int
	Error        string
	Completed    bool
}

// VideoOptions 定义视频生成的可选参数
type VideoOptions struct {
	Model              string
	ImageURL           string   // 单图参考
	Duration           int      // 视频时长（秒）
	FPS                int      // 帧率
	Resolution         string   // 分辨率 (如 "1080P", "720P", "1024x1024")
	AspectRatio        string   // 宽高比 (如 "16:9", "9:16")
	Style              string   // 风格
	MotionLevel        int      // 运动幅度 (0-100)
	CameraMotion       string   // 运镜类型
	Seed               int64    // 随机种子
	FirstFrameURL      string   // 首帧图片URL
	LastFrameURL       string   // 尾帧图片URL
	ReferenceImageURLs []string // 多图参考列表
}

// VideoOption 定义配置选项的函数签名
type VideoOption func(*VideoOptions)

func WithModel(model string) VideoOption {
	return func(o *VideoOptions) {
		o.Model = model
	}
}

func WithImageURL(url string) VideoOption {
	return func(o *VideoOptions) {
		o.ImageURL = url
	}
}

func WithDuration(duration int) VideoOption {
	return func(o *VideoOptions) {
		o.Duration = duration
	}
}

func WithFPS(fps int) VideoOption {
	return func(o *VideoOptions) {
		o.FPS = fps
	}
}

func WithResolution(resolution string) VideoOption {
	return func(o *VideoOptions) {
		o.Resolution = resolution
	}
}

func WithAspectRatio(ratio string) VideoOption {
	return func(o *VideoOptions) {
		o.AspectRatio = ratio
	}
}

func WithStyle(style string) VideoOption {
	return func(o *VideoOptions) {
		o.Style = style
	}
}

func WithMotionLevel(level int) VideoOption {
	return func(o *VideoOptions) {
		o.MotionLevel = level
	}
}

func WithCameraMotion(motion string) VideoOption {
	return func(o *VideoOptions) {
		o.CameraMotion = motion
	}
}

func WithSeed(seed int64) VideoOption {
	return func(o *VideoOptions) {
		o.Seed = seed
	}
}

func WithFirstFrame(url string) VideoOption {
	return func(o *VideoOptions) {
		o.FirstFrameURL = url
	}
}

func WithLastFrame(url string) VideoOption {
	return func(o *VideoOptions) {
		o.LastFrameURL = url
	}
}

func WithReferenceImages(urls []string) VideoOption {
	return func(o *VideoOptions) {
		o.ReferenceImageURLs = urls
	}
}

// 辅助函数：创建一个默认的 HTTP Client
func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Minute, // 视频生成请求可能较长，默认5分钟
	}
}

// NewClient 辅助函数：根据 provider 实例化客户端的工厂方法
func NewClient(provider, baseURL, apiKey, model, endpoint, queryEndpoint string) (VideoClient, error) {
	switch provider {
	case "openai", "sora":
		return NewOpenAISoraClient(baseURL, apiKey, model), nil
	case "minimax", "hailuo":
		return NewMinimaxClient(baseURL, apiKey, model), nil
	case "volces", "doubao":
		return NewVolcesArkClient(baseURL, apiKey, model, endpoint, queryEndpoint), nil
	case "runway":
		return NewRunwayClient(baseURL, apiKey, model), nil
	case "pika":
		return NewPikaClient(baseURL, apiKey, model), nil
	case "getgoapi":
		return NewGetGoAPIClient(baseURL, apiKey, model, endpoint, queryEndpoint), nil
	default:
		return nil, fmt.Errorf("unsupported video provider: %s", provider)
	}
}
