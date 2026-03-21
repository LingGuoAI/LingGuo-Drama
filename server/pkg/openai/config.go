package openai

// Config AI 全局配置
type Config struct {
	Provider string // "openai", "gemini", "doubao" 或 "getgoapi"

	// OpenAI 配置
	OpenAIBaseURL    string
	OpenAIKey        string
	OpenAIModel      string // 默认文本模型
	OpenAIImageModel string // 生图专用模型字段 (如 dall-e-3)

	// GetGoAPI 配置 (用于中转渠道)
	GetGoAPIBaseURL    string
	GetGoAPIKey        string
	GetGoAPIModel      string // GetGo 文本模型 (如 gpt-4o)
	GetGoAPIImageModel string // GetGo 生图模型 (如 gpt-4o-image)

	// Gemini 配置
	GeminiBaseURL string
	GeminiKey     string
	GeminiModel   string

	// 豆包 (Volcengine) 配置
	DoubaoBaseURL    string // 通常是 https://ark.cn-beijing.volces.com/api/v3
	DoubaoKey        string // 对应接入点的 API Key
	DoubaoModel      string // 对应接入点的 Endpoint ID (如 ep-2024xxxx-xxx)
	DoubaoImageModel string

	// Vertex AI (Google Cloud) 配置
	VertexKey        string // API 密钥
	VertexModel      string // 如 gemini-2.5-flash
	VertexImageModel string // 如 imagen-3.0-generate-001
}
