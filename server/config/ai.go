package config

import "spiritFruit/pkg/config"

func init() {
	config.Add("ai", func() map[string]interface{} {
		return map[string]interface{}{
			// 读取 AI_PROVIDER，默认为 openai
			"provider": config.Env("AI_PROVIDER", "openai"),

			"openai": map[string]interface{}{
				"base_url": config.Env("OPENAI_BASE_URL", "https://api.openai.com/v1"),
				"api_key":  config.Env("OPENAI_API_KEY", ""),
				"model":    config.Env("OPENAI_MODEL", "gpt-3.5-turbo"),
			},

			"gemini": map[string]interface{}{
				"base_url": config.Env("GEMINI_BASE_URL", "https://generativelanguage.googleapis.com/v1beta"),
				"api_key":  config.Env("GEMINI_API_KEY", ""),
				"model":    config.Env("GEMINI_MODEL", "gemini-pro"),
			},
		}
	})
}
