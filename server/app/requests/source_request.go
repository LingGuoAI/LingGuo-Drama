package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// SourceRequest 素材创建/更新请求
type SourceRequest struct {
	ProjectId  uint64 `json:"projectId" valid:"projectId"`   // 项目id
	ScriptId   uint64 `json:"scriptId" valid:"scriptId"`     // 剧本(集数)id
	ShotId     uint64 `json:"shotId" valid:"shotId"`         // 分镜ID
	ShotNumber uint64 `json:"shotNumber" valid:"shotNumber"` // 第几个镜头
	VideoUrl   string `json:"videoUrl" valid:"videoUrl"`     // 视频路径
}

// SourceSave 验证素材保存/更新的规则
func SourceSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"projectId":  []string{"required", "numeric"},
		"scriptId":   []string{"required", "numeric"},
		"shotId":     []string{"required", "numeric"},
		"shotNumber": []string{"numeric"},
		"videoUrl":   []string{"required"},
	}

	messages := govalidator.MapData{
		"projectId.required": []string{"项目ID不能为空"},
		"projectId.numeric":  []string{"项目ID必须是数字"},
		"scriptId.required":  []string{"剧本ID不能为空"},
		"scriptId.numeric":   []string{"剧本ID必须是数字"},
		"shotId.required":    []string{"分镜ID不能为空"},
		"shotId.numeric":     []string{"分镜ID必须是数字"},
		"shotNumber.numeric": []string{"镜头序号必须是数字"},
		"videoUrl.required":  []string{"视频路径不能为空"},
	}

	return validate(data, rules, messages)
}
