package v1

import (
	"github.com/gin-gonic/gin"

	"spiritFruit/app/models/source"
	"spiritFruit/app/requests"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
)

type SourceController struct {
	BaseADMINController
}

// Store 创建素材
// @Summary 创建素材
// @Description 创建新的素材(视频)
// @Tags sources
// @Accept json
// @Produce json
// @Param request body requests.SourceRequest true "素材信息"
// @Success 201 {object} response.Response{data=source.Source} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sources [post]
func (ctrl *SourceController) Store(c *gin.Context) {
	request := requests.SourceRequest{}
	// 使用 SourceSave 规则进行验证
	if ok := requests.Validate(c, &request, requests.SourceSave); !ok {
		return
	}

	// 因为模型定义的字段是指针类型（*uint64, *string），这里需要取变量地址
	projectId := request.ProjectId
	scriptId := request.ScriptId // 🔴 修改为 ScriptId
	shotId := request.ShotId
	shotNumber := request.ShotNumber
	videoUrl := request.VideoUrl

	sourceModel := source.Source{
		ProjectId:  &projectId,
		ScriptId:   &scriptId, // 🔴 修改为 ScriptId
		ShotId:     &shotId,
		ShotNumber: &shotNumber,
		VideoUrl:   &videoUrl,
	}

	// 执行创建
	sourceModel.Create()

	if sourceModel.ID > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    sourceModel,
			"message": "success",
		})
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

// Delete 删除素材
// @Summary 删除素材
// @Description 删除指定的素材
// @Tags sources
// @Accept json
// @Produce json
// @Param id path string true "sources ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "素材不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sources/{id} [delete]
func (ctrl *SourceController) Delete(c *gin.Context) {
	// 查询要删除的记录
	var sourceModel source.Source
	database.DB.Where("id = ?", c.Param("id")).First(&sourceModel)

	if sourceModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	// 调用模型里的 Delete 方法
	rowsAffected := sourceModel.Delete()
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
