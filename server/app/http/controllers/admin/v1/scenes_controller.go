package v1

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/app/models"
	"spiritFruit/app/models/projects"
	"spiritFruit/app/models/scenes"
	"spiritFruit/app/requests"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
	"strconv"
	"strings"
	"time"
)

type ScenesController struct {
	BaseADMINController
}

// Index 场景列表
// @Summary 场景列表
// @Description 获取场景列表，支持多种搜索条件
// @Tags scenes
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param projectId query string false "所属项目ID"
// @Param name query string false "场景名"
// @Param location query string false "地点"
// @Param time query string false "时间"
// @Success 200 {object} response.Response{data=[]scenes.Scenes} "场景列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scenes [get]
func (ctrl *ScenesController) Index(c *gin.Context) {
	// 构建搜索条件
	where := ctrl.buildSearchConditions(c)

	// 获取分页参数
	perPage := 10
	if perPageStr := c.Query("pageSize"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	data, pager := scenes.Paginate(c, perPage, where)
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
func (ctrl *ScenesController) buildSearchConditions(c *gin.Context) map[string]interface{} {
	where := map[string]interface{}{}

	// 所属项目ID搜索
	if projectId := strings.TrimSpace(c.Query("projectId")); projectId != "" {
		where["project_id"] = projectId
	}

	// 场景名搜索 (模糊查询)
	if name := strings.TrimSpace(c.Query("name")); name != "" {
		where["name"] = []interface{}{"like", "%" + name + "%"}
	}

	// 地点搜索
	if location := strings.TrimSpace(c.Query("location")); location != "" {
		where["location"] = location
	}

	// 时间搜索 (日/夜)
	if timeVal := strings.TrimSpace(c.Query("time")); timeVal != "" {
		where["time"] = timeVal
	}

	return where
}

// GetProjectsSelectList 获取短剧项目选择列表
// @Summary 获取短剧项目选择列表
// @Description 获取短剧项目的简化列表，用于场景中的短剧项目选择
// @Tags scenes
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]projects.Projects} "短剧项目选择列表"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scenes/getProjectsSelectList [get]
func (ctrl *ScenesController) GetProjectsSelectList(c *gin.Context) {
	list := projects.All()
	response.JSON(c, gin.H{
		"code":    0,
		"data":    list,
		"message": "success",
	})
}

// Show 场景详情
// @Summary 场景详情
// @Description 获取场景详情
// @Tags scenes
// @Accept json
// @Produce json
// @Param id path string true "scenes ID"
// @Success 200 {object} response.Response{data=scenes.Scenes} "场景详情"
// @Failure 404 {object} response.ErrorResponse "场景不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scenes/{id} [get]
func (ctrl *ScenesController) Show(c *gin.Context) {
	scenesModel := scenes.Get(c.Param("id"))
	if scenesModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}
	response.JSON(c, gin.H{
		"code":    0,
		"data":    scenesModel,
		"message": "success",
	})
}

// Store 创建场景
// @Summary 创建场景
// @Description 创建新的场景
// @Tags scenes
// @Accept json
// @Produce json
// @Param request body requests.ScenesRequest true "场景信息"
// @Success 201 {object} response.Response{data=scenes.Scenes} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scenes [post]
func (ctrl *ScenesController) Store(c *gin.Context) {
	request := requests.ScenesRequest{}
	if ok := requests.Validate(c, &request, requests.ScenesSave); !ok {
		return
	}

	// 指针类型转换
	pid := request.ProjectId
	status := request.Status

	scenesModel := scenes.Scenes{
		ProjectId:    &pid,
		Name:         &request.Name,
		Location:     &request.Location,
		Time:         &request.Time,
		Atmosphere:   &request.Atmosphere,
		VisualPrompt: &request.VisualPrompt,
		Status:       &status,
	}

	scenesModel.Create()
	if scenesModel.ID > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    scenesModel,
			"message": "success",
		})
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

// Update 更新场景
// @Summary 更新场景
// @Description 更新场景信息
// @Tags scenes
// @Accept json
// @Produce json
// @Param id path string true "scenes ID"
// @Param request body requests.ScenesRequest true "场景信息"
// @Success 200 {object} response.Response{data=scenes.Scenes} "更新成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "场景不存在"
// @Failure 422 {object} response.ValidationErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scenes/{id} [put]
func (ctrl *ScenesController) Update(c *gin.Context) {
	// 验证数据是否存在
	id := c.Param("id")
	existingscenes := scenes.Get(id)
	if existingscenes.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	request := requests.ScenesRequest{}
	if bindOk := requests.Validate(c, &request, requests.ScenesSave); !bindOk {
		return
	}

	// 使用新的模型实例进行更新，避免关联对象的影响
	updatescenes := &scenes.Scenes{
		BaseModel: models.BaseModel{ID: existingscenes.ID},
	}

	// 赋值字段
	pid := request.ProjectId
	status := request.Status

	updatescenes.ProjectId = &pid
	updatescenes.Name = &request.Name
	updatescenes.Location = &request.Location
	updatescenes.Time = &request.Time
	updatescenes.Atmosphere = &request.Atmosphere
	updatescenes.VisualPrompt = &request.VisualPrompt
	updatescenes.Status = &status
	updatescenes.UpdatedAt = time.Now()
	updatescenes.CreatedAt = existingscenes.CreatedAt // 保持创建时间不变

	// 执行更新
	result := database.DB.Save(updatescenes)

	if result.Error != nil {
		response.Abort500(c, "更新失败："+result.Error.Error())
		return
	}

	if result.RowsAffected > 0 {
		// 重新获取更新后的完整数据（包括关联）
		updatedscenes := scenes.Get(id)
		response.JSON(c, gin.H{
			"code":    0,
			"data":    updatedscenes,
			"message": "success",
		})
	} else {
		// 如果没有字段变更，Save 也可能返回 RowsAffected=0，但这通常不算失败
		// 这里可以再次查询返回最新数据
		updatedscenes := scenes.Get(id)
		response.JSON(c, gin.H{
			"code":    0,
			"data":    updatedscenes,
			"message": "success", // 这里视为成功，因为目标状态已一致
		})
	}
}

// Delete 删除场景
// @Summary 删除场景
// @Description 删除场景
// @Tags scenes
// @Accept json
// @Produce json
// @Param id path string true "scenes ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "场景不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scenes/{id} [delete]
func (ctrl *ScenesController) Delete(c *gin.Context) {
	scenesModel := scenes.Get(c.Param("id"))
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
