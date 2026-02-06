package openai

// Config AI 全局配置
type Config struct {
	Provider string // "openai" 或 "gemini"

	// OpenAI 配置
	OpenAIBaseURL string
	OpenAIKey     string
	OpenAIModel   string

	// Gemini 配置
	GeminiBaseURL string // 通常是 https://generativelanguage.googleapis.com/v1beta
	GeminiKey     string
	GeminiModel   string
}
