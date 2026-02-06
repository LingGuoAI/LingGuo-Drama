package openai

import (
	"fmt"
	"net/http"
	"strings"
)

type GeminiClient struct {
	Config Config
	client *http.Client
}

// Gemini 特有的请求结构
type geminiContent struct {
	Parts []geminiPart `json:"parts"`
	Role  string       `json:"role,omitempty"`
}
type geminiPart struct {
	Text string `json:"text"`
}
type geminiGenerateReq struct {
	Contents []geminiContent `json:"contents"`
}

// Gemini 特有的响应结构
type geminiGenerateResp struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (c *GeminiClient) GenerateScript(req ScriptRequest) (string, error) {
	// 1. 转换通用请求 -> Gemini 专用请求
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

	// 2. 构造 URL (Gemini 使用 Query 参数传 Key)
	// 默认 Model: "gemini-pro"
	model := c.Config.GeminiModel
	if model == "" {
		model = "gemini-pro"
	}
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s",
		c.Config.GeminiBaseURL, model, c.Config.GeminiKey)

	// 3. 发送请求
	resp, err := doRequest[*geminiGenerateResp](c.client, "POST", url, nil, payload)
	if err != nil {
		return "", err
	}

	// 4. 解析结果
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return resp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("empty response from gemini")
}

// GenerateImage 实现 GenerateImage (使用 Imagen 3 模型)
func (c *GeminiClient) GenerateImage(req ImageRequest) ([]string, error) {
	// 1. 设置模型
	// Google 生图通常使用 "imagegeneration@006" 或 "imagen-3.0-generate-001"
	// 如果 Config 未指定，默认使用 imagen-3.0
	model := c.Config.GeminiModel
	if model == "" || strings.Contains(model, "gemini") {
		model = "imagen-3.0-generate-001"
	}

	// 2. 构造请求体
	// Imagen API 结构比较特殊，遵循 Vertex AI 的 instances/parameters 模式
	n := req.N
	if n <= 0 {
		n = 1
	}
	if n > 4 {
		n = 4
	} // Imagen 单次限制通常为 4

	payload := geminiImageReq{
		Instances: []geminiImageInstance{
			{Prompt: req.Prompt},
		},
		Parameters: geminiImageParams{
			SampleCount: n,
		},
	}

	// 3. 构造 URL
	// 注意：Gemini/Vertex AI 的 BaseURL 结构与 OpenAI 差异较大
	// 标准格式: https://{region}-aiplatform.googleapis.com/v1/projects/{project}/locations/{location}/publishers/google/models/{model}:predict
	// 这里假设 c.Config.GeminiBaseURL 已经配置为类似 "https://generativelanguage.googleapis.com/v1beta" 的简化端点
	// 或者我们尝试使用 Google AI Studio 的简化路径 (如果可用)
	// ⚠️ 实际生产中，Google 生图目前主要通过 Vertex AI (GCP) 提供，而非简单的 API Key 调用。
	// 如果您使用的是 API Key 方式 (Google AI Studio)，目前生图功能可能受限或处于 Beta 阶段。
	// 下面URL是基于 Google AI Studio 风格的尝试：
	url := fmt.Sprintf("%s/models/%s:predict?key=%s",
		c.Config.GeminiBaseURL, model, c.Config.GeminiKey)

	// 4. 发送请求
	resp, err := doRequest[*geminiImageResp](c.client, "POST", url, nil, payload)
	if err != nil {
		return nil, err
	}

	// 5. 处理错误响应
	if resp.Error.Code != 0 {
		return nil, fmt.Errorf("gemini api error: %s", resp.Error.Message)
	}

	// 6. 提取图片并转换为 Data URI
	var urls []string
	for _, pred := range resp.Predictions {
		if pred.BytesBase64Encoded != "" {
			// 拼接成前端可直接使用的格式
			mime := pred.MimeType
			if mime == "" {
				mime = "image/png"
			}
			dataURI := fmt.Sprintf("data:%s;base64,%s", mime, pred.BytesBase64Encoded)
			urls = append(urls, dataURI)
		}
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("gemini returned no images")
	}

	return urls, nil
}
