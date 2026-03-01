package v1

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/app/models/shot_frame_image"
	"spiritFruit/app/requests"
	"spiritFruit/pkg/response"
)

type ShotFrameImagesController struct {
	BaseADMINController
}

// Store 创建分镜图片
// @Summary 创建分镜图片
// @Description 创建新的分镜图片
// @Tags shot_frame_images
// @Accept json
// @Produce json
// @Param request body requests.ShotFrameImagesRequest true "分镜图片信息"
// @Success 201 {object} response.Response{data=shot_frame_image.ShotFrameImages} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shot_frame_images [post]
func (ctrl *ShotFrameImagesController) Store(c *gin.Context) {
	request := requests.ShotFrameImagesRequest{}
	if ok := requests.Validate(c, &request, requests.ShotFrameImagesSave); !ok {
		return
	}

	projectId := request.ProjectId
	shotId := request.ShotId
	frameType := request.FrameType
	imageUrl := request.ImageUrl
	imageType := request.ImageType

	imageModel := shot_frame_image.ShotFrameImages{
		ProjectId: &projectId,
		ShotId:    &shotId,
		FrameType: &frameType,
		ImageType: &imageType,
		ImageUrl:  &imageUrl,
	}

	// 执行创建
	imageModel.Create()

	if imageModel.ID > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    imageModel,
			"message": "success",
		})
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

// Delete 删除分镜图片
// @Summary 删除分镜图片
// @Description 删除分镜图片
// @Tags scenes
// @Accept json
// @Produce json
// @Param id path string true "scenes ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "分镜图片不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scenes/{id} [delete]
func (ctrl *ShotFrameImagesController) Delete(c *gin.Context) {
	scenesModel := shot_frame_image.Get(c.Param("id"))
	if scenesModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	rowsAffected := scenesModel.Delete()
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
