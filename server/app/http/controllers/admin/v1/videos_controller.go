// app/http/controllers/admin/v1/video_controller.go
package v1

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/app/services"
	"spiritFruit/pkg/response"
)

type VideoController struct{}

// FinalizeEpisode 触发视频合成
func (ctrl *VideoController) FinalizeEpisode(c *gin.Context) {
	var req services.FinalizeEpisodeReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数解析失败: "+err.Error())
		return
	}

	videoService := new(services.VideoService)
	result, err := videoService.FinalizeEpisode(req)
	if err != nil {
		response.Abort500(c, "合成任务创建失败: "+err.Error())
		return
	}

	response.JSON(c, result)
}
