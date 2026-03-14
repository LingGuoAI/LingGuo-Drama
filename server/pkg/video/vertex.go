package video

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"spiritFruit/pkg/config"
	"spiritFruit/pkg/upload"

	"cloud.google.com/go/storage"
	"google.golang.org/genai"
)

type VertexVideoClient struct {
	Model string
}

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

// 初始化官方 GenAI 客户端
func (c *VertexVideoClient) getGenaiClient(ctx context.Context) (*genai.Client, error) {
	projectID := config.GetString("ai.vertex.project_id")
	region := config.GetString("ai.vertex.region", "us-central1")

	if projectID == "" {
		return nil, fmt.Errorf("VERTEX_PROJECT_ID is missing in config")
	}

	clientConfig := &genai.ClientConfig{
		Project:     projectID,
		Location:    region,
		Backend:     genai.BackendVertexAI,
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	}

	return genai.NewClient(ctx, clientConfig)
}

// 🔴 通过 ADC 权限从 GCS 存储桶下载文件
func downloadFromGCS(ctx context.Context, gsURI string) ([]byte, error) {
	// 解析 gs://bucket/object
	trimmed := strings.TrimPrefix(gsURI, "gs://")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid gs URI format: %s", gsURI)
	}
	bucketName := parts[0]
	objectName := parts[1]

	// 创建 Storage 客户端 (它会自动使用你配置好的 ADC 凭证)
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	// 读取对象流
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read GCS object: %w", err)
	}
	defer rc.Close()

	return io.ReadAll(rc)
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

	genConfig := &genai.GenerateVideosConfig{
		AspectRatio:  options.AspectRatio,
		OutputGCSURI: outputURI,
	}

	fmt.Printf("Vertex SDK: Submitting GenerateVideos task to model %s...\n", model)

	// 文本生成视频 (Image 传 nil)
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

	if op.Error != nil {
		if msg, ok := op.Error["message"].(string); ok && msg != "" {
			videoResult.Error = msg
			videoResult.Status = "failed"
			return videoResult, nil
		}
	}

	// 🔴 如果完成，下载 GCS 视频并转存到本地
	if op.Done {
		videoResult.Status = "completed"
		if op.Response != nil && len(op.Response.GeneratedVideos) > 0 {
			gsURI := op.Response.GeneratedVideos[0].Video.URI

			// 如果是谷歌云存储桶链接
			if strings.HasPrefix(gsURI, "gs://") {
				videoBytes, dlErr := downloadFromGCS(ctx, gsURI)
				if dlErr == nil {
					// 复用你写好的本地存储逻辑，生成 uploads/videos/xxxx.mp4
					savedPath, saveErr := upload.SaveFileDirByte(videoBytes, "videos", ".mp4")
					if saveErr == nil {
						if !strings.HasPrefix(savedPath, "/") {
							savedPath = "/" + savedPath
						}
						// 替换为本地标准路径
						videoResult.VideoURL = savedPath
					} else {
						videoResult.Error = "save local video err: " + saveErr.Error()
						videoResult.Status = "failed"
					}
				} else {
					videoResult.Error = "download GCS err: " + dlErr.Error()
					videoResult.Status = "failed"
				}
			} else {
				// 兜底：如果直接返回了 http 链接，则原样返回交由 Job 处理
				videoResult.VideoURL = gsURI
			}
		}
	}

	return videoResult, nil
}
