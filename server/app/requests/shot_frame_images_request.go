package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

// ShotFrameImagesRequest 分镜帧图片创建/更新请求
type ShotFrameImagesRequest struct {
	ProjectId uint64 `json:"projectId" valid:"projectId"` // 所属项目ID
	ShotId    uint64 `json:"shotId" valid:"shotId"`       // 关联的分镜ID
	FrameType string `json:"frameType" valid:"frameType"` // 帧类型：first-首帧 last-尾帧 action-动作序列 key-关键帧 single-单参考图
	ImageType string `json:"imageType" valid:"imageType"` // 图片类型：shot-镜头图 reference-参考图
	ImageUrl  string `json:"imageUrl" valid:"imageUrl"`   // 图片地址
}

// ShotFrameImagesSave 验证分镜帧图片保存/更新的规则
func ShotFrameImagesSave(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"projectId": []string{"required", "numeric"},
		"shotId":    []string{"required", "numeric"},
		"frameType": []string{"required", "in:first,last,action,key,single"},
		"imageType": []string{"required", "in:shot,reference"},
		"imageUrl":  []string{"required"},
	}

	messages := govalidator.MapData{
		"projectId.required": []string{"所属项目ID不能为空"},
		"projectId.numeric":  []string{"项目ID必须是数字"},
		"shotId.required":    []string{"分镜ID不能为空"},
		"shotId.numeric":     []string{"分镜ID必须是数字"},
		"frameType.required": []string{"帧类型不能为空"},
		"frameType.in":       []string{"帧类型必须是 first(首帧)、last(尾帧)、action(动作)、key(关键帧)、single(单图) 之一"},
		"imageType.required": []string{"图片类型不能为空"},
		"imageType.in":       []string{"图片类型必须是 shot(镜头图) 或 reference(参考图)"},
		"imageUrl.required":  []string{"图片地址不能为空"},
	}

	return validate(data, rules, messages)
}
