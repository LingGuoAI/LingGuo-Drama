package v1

import (
	"github.com/gin-gonic/gin"

	"spiritFruit/app/models/shot_generate_video"
	"spiritFruit/app/requests"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
)

type ShotGenerateVideosController struct {
	BaseADMINController
}

// Store 创建分镜生成视频记录
// @Summary 创建分镜生成视频
// @Description 创建新的分镜生成视频记录
// @Tags shot_generate_videos
// @Accept json
// @Produce json
// @Param request body requests.ShotGenerateVideoRequest true "分镜视频信息"
// @Success 201 {object} response.Response{data=shot_generate_video.ShotGenerateVideo} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shot_generate_videos [post]
func (ctrl *ShotGenerateVideosController) Store(c *gin.Context) {
	request := requests.ShotGenerateVideoRequest{}
	// 使用请求验证规则
	if ok := requests.Validate(c, &request, requests.ShotGenerateVideoSave); !ok {
		return
	}

	// 提取变量取址，因为模型中定义的是指针类型 (*uint64, *string)
	projectId := request.ProjectId
	scriptId := request.ScriptId
	shotId := request.ShotId
	videoUrl := request.VideoUrl

	videoModel := shot_generate_video.ShotGenerateVideo{
		ProjectId: &projectId,
		ScriptId:  &scriptId,
		ShotId:    &shotId,
		VideoUrl:  &videoUrl,
	}

	// 执行创建
	videoModel.Create()

	if videoModel.ID > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    videoModel,
			"message": "success",
		})
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

// Delete 删除分镜生成视频记录
// @Summary 删除分镜生成视频
// @Description 删除分镜生成视频记录
// @Tags shot_generate_videos
// @Accept json
// @Produce json
// @Param id path string true "记录 ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "记录不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shot_generate_videos/{id} [delete]
func (ctrl *ShotGenerateVideosController) Delete(c *gin.Context) {
	var videoModel shot_generate_video.ShotGenerateVideo

	// 查询要删除的记录
	database.DB.Where("id = ?", c.Param("id")).First(&videoModel)

	if videoModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	// 执行删除
	rowsAffected := videoModel.Delete()
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
