package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"

	"spiritFruit/app/models/shot_video_merge"
	"spiritFruit/app/requests"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
)

type ShotVideoMergesController struct {
	BaseADMINController
}

// Index 合并视频列表
// @Summary 合并视频列表
// @Description 获取合并视频列表，支持多种搜索条件
// @Tags scenes
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param scriptId query string false "所属剧集"
// @Success 200 {object} response.Response{data=[]scenes.Scenes} "合并视频列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scenes [get]
func (ctrl *ShotVideoMergesController) Index(c *gin.Context) {
	// 构建搜索条件
	where := ctrl.buildSearchConditions(c)

	// 获取分页参数
	perPage := 10
	if perPageStr := c.Query("pageSize"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	data, pager := shot_video_merge.Paginate(c, perPage, where)
	response.JSON(c, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"total": pager.TotalCount,
			"list":  data,
		},
		"message": "success",
	})
}

// buildSearchConditions 构建搜索条件
func (ctrl *ShotVideoMergesController) buildSearchConditions(c *gin.Context) map[string]interface{} {
	where := map[string]interface{}{}

	// 所属项目ID搜索
	if scriptId := strings.TrimSpace(c.Query("scriptId")); scriptId != "" {
		where["script_id"] = scriptId
	}

	return where
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
