package openai

import (
	"fmt"
	"net/http"
	"strings"
)

type GetGoAPIClient struct {
	Config Config
	client *http.Client
}

// GetGoAPI 的响应格式与 OpenAI 一致
type getGoChatResp struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// GenerateScript 实现文本/剧本生成
func (c *GetGoAPIClient) GenerateScript(req ScriptRequest) (string, error) {
	// 优先使用 GetGo 专有模型配置
	model := c.Config.GetGoAPIModel
	if model == "" {
		model = c.Config.OpenAIModel // 兜底使用通用模型名
	}

	payload := map[string]interface{}{
		"model":      model,
		"messages":   req.Messages,
		"max_tokens": req.MaxTokens,
	}

	// 使用 GetGo 专有 Key
	apiKey := c.Config.GetGoAPIKey
	if apiKey == "" {
		apiKey = c.Config.OpenAIKey
	}

	headers := map[string]string{
		"Authorization": "Bearer " + apiKey,
	}

	// 使用 GetGo 专有 BaseURL
	baseURL := c.Config.GetGoAPIBaseURL
	if baseURL == "" {
		baseURL = c.Config.OpenAIBaseURL
	}

	url := strings.TrimRight(baseURL, "/") + "/chat/completions"

	resp, err := doRequest[*getGoChatResp](c.client, "POST", url, headers, payload)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("getgoapi returned no choices")
}

// GenerateImage 实现 GenerateImage (支持 DALL-E 3/2)
func (c *GetGoAPIClient) GenerateImage(req ImageRequest) ([]string, error) {
	// 1. 设置默认值与构建请求体
	model := c.Config.GetGoAPIImageModel
	if model == "" {
		model = "dall-e-3" // 默认使用 DALL-E 3
	}
	fmt.Println("model", model)
	// DALL-E 3 的数量限制为 1
	n := req.N
	if n <= 0 {
		n = 1
	}

	// 默认尺寸
	size := req.Size
	if size == "" {
		size = "1024x1024"
	}

	payload := openAIImageReq{
		Model:  model,
		Prompt: req.Prompt,
		N:      n,
		Size:   size,
	}

	// 2. 设置 Headers
	headers := map[string]string{
		"Authorization": "Bearer " + c.Config.GetGoAPIKey,
	}

	// 3. 拼接 URL
	url := c.Config.GetGoAPIBaseURL + "/images/generations"

	// 4. 发送请求 (复用泛型 doRequest)
	resp, err := doRequest[*openAIImageResp](c.client, "POST", url, headers, payload)
	if err != nil {
		return nil, err
	}

	// 5. 提取图片 URL
	var urls []string
	for _, item := range resp.Data {
		urls = append(urls, item.URL)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("openai returned empty image list")
	}

	return urls, nil
}
