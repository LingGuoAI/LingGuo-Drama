package openai

// ================== 通用结构 ==================

// CommonErrorResponse 通用错误响应
type CommonErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

// ================== 文本生成 (LLM) ==================

// ================== 统一的入参和出参 ==================

// ScriptRequest 剧本生成请求（业务层通用）
type ScriptRequest struct {
	Messages    []ChatMessage
	Temperature float64
}

// ImageRequest 图片生成请求（业务层通用）
type ImageRequest struct {
	Prompt string
	N      int
	Size   string
}

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openAIImageReq OpenAI 生图接口请求参数
type openAIImageReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"` // 生成数量 (DALL-E 3 只能是 1)
	Size   string `json:"size"`
	// Quality string `json:"quality,omitempty"` // 可选: standard 或 hd
	// Style   string `json:"style,omitempty"`   // 可选: vivid 或 natural
}

// openAIImageResp OpenAI 生图接口响应结构
type openAIImageResp struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL string `json:"url"` // 图片链接
		// B64JSON string `json:"b64_json"` // 如果需要 Base64
	} `json:"data"`
}

// geminiImageReq Imagen 3 生图请求参数
// 官方文档参考: https://cloud.google.com/vertex-ai/generative-ai/docs/image/img-gen-prompt-api
type geminiImageReq struct {
	Instances  []geminiImageInstance `json:"instances"`
	Parameters geminiImageParams     `json:"parameters"`
}

type geminiImageInstance struct {
	Prompt string `json:"prompt"`
}

type geminiImageParams struct {
	SampleCount int `json:"sampleCount,omitempty"` // 生成数量
	// AspectRatio string `json:"aspectRatio,omitempty"` // 宽高比: "1:1", "16:9", etc.
	// PersonGeneration string `json:"personGeneration,omitempty"` // "allow_adult", "dont_allow"
}

// geminiImageResp Imagen 3 生图响应结构
type geminiImageResp struct {
	Predictions []struct {
		BytesBase64Encoded string `json:"bytesBase64Encoded"` // 图片是以 Base64 返回的
		MimeType           string `json:"mimeType"`
	} `json:"predictions"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
