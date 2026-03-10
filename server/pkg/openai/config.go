package openai

// Config AI 全局配置
type Config struct {
	Provider string // "openai", "gemini" 或 "doubao"

	// OpenAI 配置
	OpenAIBaseURL string
	OpenAIKey     string
	OpenAIModel   string

	// Gemini 配置
	GeminiBaseURL string
	GeminiKey     string
	GeminiModel   string

	// 🔴 豆包 (Volcengine) 配置
	DoubaoBaseURL    string // 通常是 https://ark.cn-beijing.volces.com/api/v3
	DoubaoKey        string // 对应接入点的 API Key
	DoubaoModel      string // 对应接入点的 Endpoint ID (如 ep-2024xxxx-xxx)
	DoubaoImageModel string
}
