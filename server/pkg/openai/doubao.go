package openai

import (
	"fmt"
	"net/http"
	"strings"
)

// DoubaoClient 豆包 API 客户端
type DoubaoClient struct {
	Config Config
	client *http.Client
}

// GenerateScript 实现 GenerateScript (文本生成)
func (c *DoubaoClient) GenerateScript(req ScriptRequest) (string, error) {
	// 火山引擎 Ark API 高度兼容 OpenAI 格式，所以直接复用 openAIChatReq 结构体
	payload := openAIChatReq{
		Model:     c.Config.DoubaoModel,
		Messages:  req.Messages,
		MaxTokens: req.MaxTokens,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + c.Config.DoubaoKey,
	}

	// 默认使用火山引擎北京区的 API 地址
	baseURL := c.Config.DoubaoBaseURL
	if baseURL == "" {
		baseURL = "https://ark.cn-beijing.volces.com/api/v3"
	}
	url := strings.TrimRight(baseURL, "/") + "/chat/completions"

	// 复用泛型 doRequest，并使用 openAIChatResp 解析
	resp, err := doRequest[*openAIChatResp](c.client, "POST", url, headers, payload)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("no choices returned from doubao")
}

// GenerateImage 实现 GenerateImage (支持豆包生图模型)
func (c *DoubaoClient) GenerateImage(req ImageRequest) ([]string, error) {
	model := c.Config.DoubaoImageModel
	if model == "" {
		return nil, fmt.Errorf("DOUBAO_IMAGE_MODEL (endpoint id for image generation) is not configured in .env")
	}

	n := req.N
	if n <= 0 {
		n = 1
	}
	fmt.Println("req-------", req.Size)
	// 1920 * 1920 = 3686400
	// 2560 * 1440 = 3686400
	size := req.Size
	switch size {
	case "16:9", "1920x1080":
		size = "2560x1440" // 满足 16:9 且达到像素阈值
	case "9:16", "1080x1920":
		size = "1440x2560" // 满足 9:16 且达到像素阈值
	case "", "1:1", "1024x1024":
		size = "1920x1920" // 默认 1:1 满足像素阈值
	default:
		// 如果传入了其他未知的极小分辨率，直接兜底强转
		size = "1920x1920"
	}

	// 豆包的文生图 API 兼容 OpenAI 格式
	payload := openAIImageReq{
		Model:  model,
		Prompt: req.Prompt,
		N:      n,
		Size:   size, // 传入自动调整后的合规尺寸
	}

	headers := map[string]string{
		"Authorization": "Bearer " + c.Config.DoubaoKey,
	}

	baseURL := c.Config.DoubaoBaseURL
	if baseURL == "" {
		baseURL = "https://ark.cn-beijing.volces.com/api/v3"
	}
	url := strings.TrimRight(baseURL, "/") + "/images/generations"

	resp, err := doRequest[*openAIImageResp](c.client, "POST", url, headers, payload)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, item := range resp.Data {
		urls = append(urls, item.URL)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("doubao returned empty image list")
	}

	return urls, nil
}
