package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// ShotVideoMergeRequest 视频合并创建/更新请求
type ShotVideoMergeRequest struct {
	ProjectId uint64 `json:"projectId" valid:"projectId"` // 项目ID
	ScriptId  uint64 `json:"scriptId" valid:"scriptId"`   // 剧本ID
	Title     string `json:"title" valid:"title"`         // 视频标题
}

// ShotVideoMergeSave 验证视频合并保存的规则
func ShotVideoMergeSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"projectId": []string{"required", "numeric"},
		"scriptId":  []string{"required", "numeric"},
	}

	messages := govalidator.MapData{
		"projectId.required": []string{"项目ID不能为空"},
		"projectId.numeric":  []string{"项目ID必须是数字"},
		"scriptId.required":  []string{"剧本ID不能为空"},
		"scriptId.numeric":   []string{"剧本ID必须是数字"},
	}

	return validate(data, rules, messages)
}
