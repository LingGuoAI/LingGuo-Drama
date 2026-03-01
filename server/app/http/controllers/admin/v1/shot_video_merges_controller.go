package v1

import (
	"github.com/gin-gonic/gin"

	"spiritFruit/app/models/shot_video_merge"
	"spiritFruit/app/requests"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
)

type ShotVideoMergesController struct {
	BaseADMINController
}

// Store 创建视频合并记录
func (ctrl *ShotVideoMergesController) Store(c *gin.Context) {
	request := requests.ShotVideoMergeRequest{}
	// 使用请求验证规则
	if ok := requests.Validate(c, &request, requests.ShotVideoMergeSave); !ok {
		return
	}

	projectId := request.ProjectId
	scriptId := request.ScriptId
	title := request.Title

	mergeModel := shot_video_merge.ShotVideoMerge{
		ProjectId: &projectId,
		ScriptId:  &scriptId,
		Title:     &title,
	}

	// 执行创建
	mergeModel.Create()

	if mergeModel.ID > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    mergeModel,
			"message": "success",
		})
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

// Delete 删除视频合并记录
func (ctrl *ShotVideoMergesController) Delete(c *gin.Context) {
	var mergeModel shot_video_merge.ShotVideoMerge

	// 查询要删除的记录
	database.DB.Where("id = ?", c.Param("id")).First(&mergeModel)

	if mergeModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	// 执行删除
	rowsAffected := mergeModel.Delete()
	if rowsAffected > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    "",
			"message": "success",
		})
		return
	}

	response.Abort500(c, "删除失败，请稍后尝试~")
}
