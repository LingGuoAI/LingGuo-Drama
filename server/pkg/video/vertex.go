package video

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"spiritFruit/pkg/config"

	"google.golang.org/genai"
)

type VertexVideoClient struct {
	APIKey string // 必须是 gcloud 产生的 OAuth Token (ya29. 开头)
	Model  string
}

func NewVertexVideoClient(baseURL, apiKey, model string) *VertexVideoClient {
	if model == "" {
		model = "veo-3.1-fast-generate-001"
	}
	return &VertexVideoClient{
		APIKey: apiKey,
		Model:  model,
	}
}

// 初始化官方 GenAI 客户端 (严格对齐文档的 v1 版本)
func (c *VertexVideoClient) getGenaiClient(ctx context.Context) (*genai.Client, error) {
	projectID := config.GetString("ai.vertex.project_id")
	region := config.GetString("ai.vertex.region", "us-central1")

	if projectID == "" {
		return nil, fmt.Errorf("VERTEX_PROJECT_ID is missing in .env")
	}

	// 🔴 构造自定义 Header，将凭证注入其中，绕过 SDK 的互斥校验
	headers := make(http.Header)
	if c.APIKey != "" {
		headers.Set("Authorization", "Bearer "+c.APIKey)
	}

	clientConfig := &genai.ClientConfig{
		Project:  projectID,
		Location: region,
		Backend:  genai.BackendVertexAI,
		// 🔴 配置 API 版本，并注入带有 Token 的自定义请求头
		HTTPOptions: genai.HTTPOptions{
			APIVersion: "v1",
			Headers:    headers,
		},
	}

	return genai.NewClient(ctx, clientConfig)
}
func (c *VertexVideoClient) GenerateVideo(prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration:    5,
		AspectRatio: "16:9",
	}
	for _, opt := range opts {
		opt(options)
	}

	model := c.Model
	if options.Model != "" {
		model = options.Model
	}

	// 🔴 获取存储桶地址 (Veo 模型强制要求)
	gcsBucket := config.GetString("ai.vertex.gcs_bucket")
	if gcsBucket == "" {
		return nil, fmt.Errorf("VERTEX_GCS_BUCKET is missing in .env (e.g., gs://your-bucket). Veo mandates a GCS bucket to output videos.")
	}

	// 确保加上后缀前缀，例如 gs://your-bucket/videos/
	outputURI := strings.TrimRight(gcsBucket, "/") + "/videos/"

	ctx := context.Background()
	client, err := c.getGenaiClient(ctx)
	if err != nil {
		return nil, err
	}

	// 1. 构建官方配置
	genConfig := &genai.GenerateVideosConfig{
		AspectRatio:  options.AspectRatio,
		OutputGCSURI: outputURI, // 指定输出目录
	}

	fmt.Printf("Vertex SDK: Submitting GenerateVideos task to model %s...\n", model)

	// 2. 调用生成接口。注意第四个参数 (Image) 目前传 nil，先确保纯文本生成完美编译和运行
	operation, err := client.Models.GenerateVideos(ctx, model, prompt, nil, genConfig)
	if err != nil {
		return nil, fmt.Errorf("genai sdk generate videos error: %w", err)
	}

	return &VideoResult{
		TaskID:    operation.Name,
		Status:    "processing",
		Completed: false,
		Duration:  options.Duration,
	}, nil
}

func (c *VertexVideoClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	ctx := context.Background()
	client, err := c.getGenaiClient(ctx)
	if err != nil {
		return nil, err
	}

	// 官方的 GenerateVideos 返回的是 *genai.GenerateVideosOperation
	dummyOp := &genai.GenerateVideosOperation{Name: taskID}

	// 调用视频专用的查询接口
	op, err := client.Operations.GetVideosOperation(ctx, dummyOp, nil)
	if err != nil {
		return nil, fmt.Errorf("genai sdk get videos operation error: %w", err)
	}

	videoResult := &VideoResult{
		TaskID:    op.Name,
		Status:    "processing",
		Completed: op.Done,
	}

	// 4. 🔴  op.Error.Message 报错：op.Error 是一个 map[string]any
	if op.Error != nil {
		if msg, ok := op.Error["message"].(string); ok && msg != "" {
			videoResult.Error = msg
			videoResult.Status = "failed"
			return videoResult, nil
		}
	}

	// 5. 🔴  Bytes 报错：读取视频的 gs:// 地址
	if op.Done {
		videoResult.Status = "completed"
		if op.Response != nil && len(op.Response.GeneratedVideos) > 0 {
			// 根据官方文档，这里返回的是 GCS URI (例如: gs://your-bucket/videos/xxx.mp4)
			videoResult.VideoURL = op.Response.GeneratedVideos[0].Video.URI
		}
	}

	return videoResult, nil
}
