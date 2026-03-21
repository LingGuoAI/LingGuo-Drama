package config

import "spiritFruit/pkg/config"

func init() {
	config.Add("ai", func() map[string]interface{} {
		return map[string]interface{}{
			// --- 默认文本大模型配置 ---
			"provider": config.Env("AI_PROVIDER", "openai"),

			"openai": map[string]interface{}{
				"base_url":    config.Env("OPENAI_BASE_URL", "https://api.openai.com/v1"),
				"api_key":     config.Env("OPENAI_API_KEY", ""),
				"model":       config.Env("OPENAI_MODEL", "gpt-3.5-turbo"),
				"image_model": config.Env("OPENAI_IMAGE_MODEL", "dall-e-3"),
			},

			"getgoapi": map[string]interface{}{
				"base_url":    config.Env("GETGOAPI_BASE_URL", "https://api.lingguoai.com/v1"),
				"api_key":     config.Env("GETGOAPI_API_KEY", ""),
				"model":       config.Env("GETGOAPI_MODEL", "gpt-4o"),
				"image_model": config.Env("GETGOAPI_IMAGE_MODEL", "gpt-4o-image"),
			},

			"doubao": map[string]interface{}{
				"base_url":    config.Env("DOUBAO_BASE_URL", "https://ark.cn-beijing.volces.com/api/v3"),
				"api_key":     config.Env("DOUBAO_API_KEY", ""),
				"model":       config.Env("DOUBAO_MODEL", ""),       // 文本模型 Endpoint
				"image_model": config.Env("DOUBAO_IMAGE_MODEL", ""), // 生图模型 Endpoint
			},

			// Vertex AI 配置 (API Key 模式)
			"vertex": map[string]interface{}{
				"api_key":     config.Env("VERTEX_API_KEY", ""),
				"model":       config.Env("VERTEX_MODEL", "gemini-1.5-pro"),
				"image_model": config.Env("VERTEX_IMAGE_MODEL", "imagen-3.0-generate-001"),
				"video_model": config.Env("VERTEX_VIDEO_MODEL", "veo-2.0-generate-001"),
				"project_id":  config.Env("VERTEX_PROJECT_ID", ""),
				"region":      config.Env("VERTEX_REGION", "us-central1"),
				"gcs_bucket":  config.Env("VERTEX_GCS_BUCKET", ""),
			},

			"gemini": map[string]interface{}{
				"base_url": config.Env("GEMINI_BASE_URL", "https://generativelanguage.googleapis.com/v1beta"),
				"api_key":  config.Env("GEMINI_API_KEY", ""),
				"model":    config.Env("GEMINI_MODEL", "gemini-pro"),
			},

			// --- 视频生成大模型配置 ---
			"video_provider": config.Env("VIDEO_PROVIDER", "getgoapi"),

			"volces": map[string]interface{}{
				"base_url": config.Env("VOLCES_BASE_URL", "https://ark.cn-beijing.volces.com/api"),
				"api_key":  config.Env("VOLCES_API_KEY", ""),
			},

			"minimax": map[string]interface{}{
				"base_url": config.Env("MINIMAX_BASE_URL", "https://api.minimax.chat/v1"),
				"api_key":  config.Env("MINIMAX_API_KEY", ""),
			},

			"runway": map[string]interface{}{
				"base_url": config.Env("RUNWAY_BASE_URL", "https://api.runwayml.com"),
				"api_key":  config.Env("RUNWAY_API_KEY", ""),
			},

			"pika": map[string]interface{}{
				"base_url": config.Env("PIKA_BASE_URL", "https://api.pika.art"),
				"api_key":  config.Env("PIKA_API_KEY", ""),
			},
		}
	})
}
