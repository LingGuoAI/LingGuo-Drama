package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// ShotGenerateVideoRequest 分镜生成视频创建/更新请求
type ShotGenerateVideoRequest struct {
	ProjectId uint64 `json:"projectId" valid:"projectId"` // 项目id
	ScriptId  uint64 `json:"scriptId" valid:"scriptId"`   // 剧本(集数)id
	ShotId    uint64 `json:"shotId" valid:"shotId"`       // 分镜ID
	VideoUrl  string `json:"videoUrl" valid:"videoUrl"`   // 视频路径
}

// ShotGenerateVideoSave 验证分镜生成视频保存/更新的规则
func ShotGenerateVideoSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"projectId": []string{"required", "numeric"},
		"scriptId":  []string{"required", "numeric"},
		"shotId":    []string{"required", "numeric"},
		"videoUrl":  []string{"required"},
	}

	messages := govalidator.MapData{
		"projectId.required": []string{"项目ID不能为空"},
		"projectId.numeric":  []string{"项目ID必须是数字"},
		"scriptId.required":  []string{"剧本ID不能为空"},
		"scriptId.numeric":   []string{"剧本ID必须是数字"},
		"shotId.required":    []string{"分镜ID不能为空"},
		"shotId.numeric":     []string{"分镜ID必须是数字"},
		"videoUrl.required":  []string{"视频路径不能为空"},
	}

	return validate(data, rules, messages)
}
