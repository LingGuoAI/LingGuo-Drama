package utils

import (
	"encoding/json"
	"strings"
)

// CleanAIJSON 清洗 AI 返回的 JSON 字符串（去除 markdown 代码块）
func CleanAIJSON(content string) string {
	content = strings.TrimSpace(content)
	// 去除 ```json 和 ```
	if strings.HasPrefix(content, "```json") {
		content = content[7:]
	} else if strings.HasPrefix(content, "```") {
		content = content[3:]
	}
	if strings.HasSuffix(content, "```") {
		content = content[:len(content)-3]
	}
	return strings.TrimSpace(content)
}

// SafeParseAIJSON 解析 AI JSON
func SafeParseAIJSON(content string, v interface{}) error {
	cleanContent := CleanAIJSON(content)
	return json.Unmarshal([]byte(cleanContent), v)
}
