package video

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"spiritFruit/pkg/config"

	"google.golang.org/genai"
)

type VertexVideoClient struct {
	Model string
}

// 注意：由于 SDK 强制分离，我们不再需要在这里保存 APIKey
func NewVertexVideoClient(baseURL, apiKey, model string) *VertexVideoClient {
	if model == "" {
		model = "veo-3.1-fast-generate-001"
	}
	return &VertexVideoClient{
		Model: model,
	}
}

// 提取并解码图片为 SDK 兼容的格式
func fetchImageBytes(imageURL string) (string, []byte, error) {
	if strings.HasPrefix(imageURL, "data:image") {
		parts := strings.SplitN(imageURL, ",", 2)
		if len(parts) != 2 {
			return "", nil, fmt.Errorf("invalid data URI")
		}
		mimeType := "image/jpeg"
		meta := strings.TrimPrefix(parts[0], "data:")
		metaParts := strings.Split(meta, ";")
		if len(metaParts) > 0 && metaParts[0] != "" {
			mimeType = metaParts[0]
		}
		data, err := base64.StdEncoding.DecodeString(parts[1])
		return mimeType, data, err
	} else if strings.HasPrefix(imageURL, "http") {
		resp, err := http.Get(imageURL)
		if err != nil {
			return "", nil, err
		}
		defer resp.Body.Close()
		mimeType := resp.Header.Get("Content-Type")
		if mimeType == "" {
			mimeType = "image/jpeg"
		}
		data, err := io.ReadAll(resp.Body)
		return mimeType, data, err
	}
	return "", nil, fmt.Errorf("unsupported image url format")
}

// 🔴 核心重构：完全遵循官方文档初始化 Vertex AI 客户端
func (c *VertexVideoClient) getGenaiClient(ctx context.Context) (*genai.Client, error) {
	projectID := config.GetString("ai.vertex.project_id")
	region := config.GetString("ai.vertex.region", "us-central1")

	if projectID == "" {
		return nil, fmt.Errorf("VERTEX_PROJECT_ID is missing in config")
	}

	// 严格按照文档：只传 Project, Location 和 Backend
	clientConfig := &genai.ClientConfig{
		Project:  projectID,
		Location: region,
		Backend:  genai.BackendVertexAI,
		// 保留 v1 API 版本的设置，因为视频接口需要
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
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

	// Veo 模型强制要求输出到 GCS 存储桶
	gcsBucket := config.GetString("ai.vertex.gcs_bucket")
	if gcsBucket == "" {
		return nil, fmt.Errorf("VERTEX_GCS_BUCKET is missing in config (e.g., gs://your-bucket)")
	}
	outputURI := strings.TrimRight(gcsBucket, "/") + "/videos/"

	ctx := context.Background()
	client, err := c.getGenaiClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}

	// 构建生成配置
	genConfig := &genai.GenerateVideosConfig{
		AspectRatio:  options.AspectRatio,
		OutputGCSURI: outputURI,
	}

	// 处理参考图（如果有）
	//var imagePart *genai.Part
	//if options.ImageURL != "" {
	//	mimeType, data, err := fetchImageBytes(options.ImageURL)
	//	if err == nil {
	//		imagePart = &genai.Part{
	//			InlineData: &genai.Blob{
	//				MIMEType: mimeType,
	//				Data:     data,
	//			},
	//		}
	//	}
	//}

	fmt.Printf("Vertex SDK: Submitting GenerateVideos task to model %s...\n", model)

	// 一键调用 SDK 提交任务
	operation, err := client.Models.GenerateVideos(ctx, model, prompt, nil, genConfig)
	if err != nil {
		return nil, fmt.Errorf("sdk generate videos error: %w", err)
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

	// 构造专用 Operation 对象并调用 GetVideosOperation
	dummyOp := &genai.GenerateVideosOperation{Name: taskID}
	op, err := client.Operations.GetVideosOperation(ctx, dummyOp, nil)
	if err != nil {
		return nil, fmt.Errorf("sdk get videos operation error: %w", err)
	}

	videoResult := &VideoResult{
		TaskID:    op.Name,
		Status:    "processing",
		Completed: op.Done,
	}

	// 解析可能存在的错误
	if op.Error != nil {
		if msg, ok := op.Error["message"].(string); ok && msg != "" {
			videoResult.Error = msg
			videoResult.Status = "failed"
			return videoResult, nil
		}
	}

	// 任务完成，返回 gs:// 格式的视频地址
	if op.Done {
		videoResult.Status = "completed"
		if op.Response != nil && len(op.Response.GeneratedVideos) > 0 {
			videoResult.VideoURL = op.Response.GeneratedVideos[0].Video.URI
		}
	}

	return videoResult, nil
}
