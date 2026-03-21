package openai

import (
	"fmt"
	"net/http"
)

// VertexClient Google Cloud Vertex AI 客户端 (基于 API Key)
type VertexClient struct {
	Config Config
	client *http.Client
}

// GenerateScript 实现文本生成
func (c *VertexClient) GenerateScript(req ScriptRequest) (string, error) {
	// 1. 转换请求
	var contents []geminiContent
	for _, msg := range req.Messages {
		role := "user"
		if msg.Role == "model" || msg.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, geminiContent{
			Role:  role,
			Parts: []geminiPart{{Text: msg.Content}},
		})
	}

	payload := geminiGenerateReq{Contents: contents}
	if req.MaxTokens > 0 {
		payload.GenerationConfig = &geminiGenerationConfig{
			MaxOutputTokens: req.MaxTokens,
		}
	}

	model := c.Config.VertexModel
	if model == "" {
		model = "gemini-1.5-pro"
	}

	apiKey := c.Config.VertexKey
	if apiKey == "" {
		return "", fmt.Errorf("VERTEX_API_KEY is required for vertex ai")
	}

	url := fmt.Sprintf("https://aiplatform.googleapis.com/v1/publishers/google/models/%s:generateContent?key=%s",
		model, apiKey)

	// 3. 发送请求 (不需要设置 Authorization: Bearer 头了)
	resp, err := doRequest[*geminiGenerateResp](c.client, "POST", url, nil, payload)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return resp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("empty response from vertex ai")
}

// GenerateImage 实现图片生成 (使用 Imagen 模型)
func (c *VertexClient) GenerateImage(req ImageRequest) ([]string, error) {
	model := c.Config.VertexImageModel
	if model == "" {
		model = "imagen-3.0-generate-001"
	}

	apiKey := c.Config.VertexKey
	if apiKey == "" {
		return nil, fmt.Errorf("VERTEX_API_KEY is required for vertex ai image generation")
	}

	n := req.N
	if n <= 0 {
		n = 1
	}
	if n > 4 {
		n = 4
	}

	// 1. 构造请求体
	payload := geminiImageReq{
		Instances: []geminiImageInstance{
			{Prompt: req.Prompt},
		},
		Parameters: geminiImageParams{
			SampleCount: n,
		},
	}

	// 2. 拼接 Imagen 的请求地址 (后缀是 :predict)
	url := fmt.Sprintf("https://aiplatform.googleapis.com/v1/publishers/google/models/%s:predict?key=%s",
		model, apiKey)

	// 3. 发送请求
	resp, err := doRequest[*geminiImageResp](c.client, "POST", url, nil, payload)
	if err != nil {
		return nil, err
	}

	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("vertex imagen api error: %s", resp.Error.Message)
	}

	// 4. 提取图片 (Base64)
	var urls []string
	for _, pred := range resp.Predictions {
		if pred.BytesBase64Encoded != "" {
			mime := pred.MimeType
			if mime == "" {
				mime = "image/png"
			}
			dataURI := fmt.Sprintf("data:%s;base64,%s", mime, pred.BytesBase64Encoded)
			urls = append(urls, dataURI)
		}
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("vertex returned no images")
	}

	return urls, nil
}
