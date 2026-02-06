package openai

// ================== 接口定义 ==================

// Provider AI 提供商接口
type Provider interface {
	// GenerateScript 生成文本/剧本
	GenerateScript(req ScriptRequest) (string, error)

	// GenerateImage 生成图片 (返回图片URL列表)
	GenerateImage(req ImageRequest) ([]string, error)
}
