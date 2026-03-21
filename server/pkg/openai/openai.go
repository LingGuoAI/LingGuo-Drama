package openai

import (
	"fmt"
	"net/http"
)

type OpenAIClient struct {
	Config Config
	client *http.Client
}

// 内部请求结构体 (OpenAI 特有)
type openAIChatReq struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens,omitempty"`
}

type openAIChatResp struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// GenerateScript 实现 GenerateScript
func (c *OpenAIClient) GenerateScript(req ScriptRequest) (string, error) {
	payload := openAIChatReq{
		Model:     c.Config.OpenAIModel,
		Messages:  req.Messages,
		MaxTokens: req.MaxTokens,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + c.Config.OpenAIKey,
	}

	url := c.Config.OpenAIBaseURL + "/chat/completions"

	resp, err := doRequest[*openAIChatResp](c.client, "POST", url, headers, payload)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no choices returned")
}

// GenerateImage 实现 GenerateImage (支持 DALL-E 3/2)
func (c *OpenAIClient) GenerateImage(req ImageRequest) ([]string, error) {
	// 1. 设置默认值与构建请求体
	model := c.Config.OpenAIImageModel
	if model == "" {
		model = "dall-e-3" // 默认使用 DALL-E 3
	}

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
		"Authorization": "Bearer " + c.Config.OpenAIKey,
	}

	// 3. 拼接 URL
	url := c.Config.OpenAIBaseURL + "/images/generations"

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
