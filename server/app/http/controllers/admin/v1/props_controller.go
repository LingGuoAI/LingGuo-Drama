package v1

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/app/models"
	"spiritFruit/app/models/props" // 引入道具模型
	"spiritFruit/app/requests"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
	"strconv"
	"strings"
	"time"
)

type PropsController struct {
	BaseADMINController
}

// Index 道具列表
// @Summary 道具列表
// @Description 获取道具列表，支持多种搜索条件
// @Tags props
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param projectId query string false "所属项目ID"
// @Param name query string false "道具名"
// @Param type query string false "道具类型"
// @Success 200 {object} response.Response{data=map[string]interface{}} "道具列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/props [get]
func (ctrl *PropsController) Index(c *gin.Context) {
	// 构建搜索条件
	where := ctrl.buildSearchConditions(c)

	// 获取分页参数
	perPage := 10
	if perPageStr := c.Query("pageSize"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	data, pager := props.Paginate(c, perPage, where)
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
func (ctrl *PropsController) buildSearchConditions(c *gin.Context) map[string]interface{} {
	where := map[string]interface{}{}

	// 所属项目ID搜索
	if projectId := strings.TrimSpace(c.Query("projectId")); projectId != "" {
		where["project_id"] = projectId
	}

	// 道具名搜索 (模糊查询)
	if name := strings.TrimSpace(c.Query("name")); name != "" {
		where["name"] = []interface{}{"like", "%" + name + "%"}
	}

	// 类型搜索
	if propType := strings.TrimSpace(c.Query("type")); propType != "" {
		where["type"] = propType
	}

	return where
}

// GetProjectsSelectList 获取项目选择列表 (这个接口通常在 ProjectController 或者通用接口里，如果放在这里，通常是获取所有道具)
// 这里假设您是想获取所有道具列表供前端选择
// @Summary 获取道具选择列表
// @Description 获取所有道具的简化列表
// @Tags props
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]props.Props} "道具列表"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/props/getSelectList [get]
func (ctrl *PropsController) GetSelectList(c *gin.Context) {
	list := props.All()
	response.JSON(c, gin.H{
		"code":    0,
		"data":    list,
		"message": "success",
	})
}

// Show 道具详情
// @Summary 道具详情
// @Description 获取道具详情
// @Tags props
// @Accept json
// @Produce json
// @Param id path string true "Props ID"
// @Success 200 {object} response.Response{data=props.Props} "道具详情"
// @Failure 404 {object} response.ErrorResponse "道具不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/props/{id} [get]
func (ctrl *PropsController) Show(c *gin.Context) {
	propsModel := props.Get(c.Param("id"))
	if propsModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}
	response.JSON(c, gin.H{
		"code":    0,
		"data":    propsModel,
		"message": "success",
	})
}

// Store 创建道具
// @Summary 创建道具
// @Description 创建新的道具
// @Tags props
// @Accept json
// @Produce json
// @Param request body requests.PropsRequest true "道具信息"
// @Success 201 {object} response.Response{data=props.Props} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/props [post]
func (ctrl *PropsController) Store(c *gin.Context) {
	request := requests.PropsRequest{}
	if ok := requests.Validate(c, &request, requests.PropsSave); !ok {
		return
	}

	// 字段赋值
	pid := request.ProjectId

	propsModel := props.Props{
		ProjectId:   &pid,
		Name:        &request.Name,
		Type:        &request.Type,
		Description: &request.Description,
		ImagePrompt: &request.ImagePrompt,
		ImageUrl:    &request.ImageUrl,
	}

	propsModel.Create()
	if propsModel.ID > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    propsModel,
			"message": "success",
		})
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

// Update 更新道具
// @Summary 更新道具
// @Description 更新道具信息
// @Tags props
// @Accept json
// @Produce json
// @Param id path string true "Props ID"
// @Param request body requests.PropsRequest true "道具信息"
// @Success 200 {object} response.Response{data=props.Props} "更新成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "道具不存在"
// @Failure 422 {object} response.ValidationErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/props/{id} [put]
func (ctrl *PropsController) Update(c *gin.Context) {
	// 验证数据是否存在
	id := c.Param("id")
	existingProp := props.Get(id)
	if existingProp.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	request := requests.PropsRequest{}
	if bindOk := requests.Validate(c, &request, requests.PropsSave); !bindOk {
		return
	}

	// 使用新的模型实例进行更新
	updateProp := &props.Props{
		BaseModel: models.BaseModel{ID: existingProp.ID},
	}

	// 赋值字段
	pid := request.ProjectId

	updateProp.ProjectId = &pid
	updateProp.Name = &request.Name
	updateProp.Type = &request.Type
	updateProp.Description = &request.Description
	updateProp.ImagePrompt = &request.ImagePrompt
	updateProp.ImageUrl = &request.ImageUrl
	updateProp.UpdatedAt = time.Now()
	updateProp.CreatedAt = existingProp.CreatedAt // 保持创建时间

	// 执行更新
	result := database.DB.Save(updateProp)

	if result.Error != nil {
		response.Abort500(c, "更新失败："+result.Error.Error())
		return
	}

	// 重新获取更新后的完整数据
	updatedProp := props.Get(id)
	response.JSON(c, gin.H{
		"code":    0,
		"data":    updatedProp,
		"message": "success",
	})
}

// Delete 删除道具
// @Summary 删除道具
// @Description 删除道具
// @Tags props
// @Accept json
// @Produce json
// @Param id path string true "Props ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "道具不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/props/{id} [delete]
func (ctrl *PropsController) Delete(c *gin.Context) {
	propsModel := props.Get(c.Param("id"))
	if propsModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	rowsAffected := propsModel.Delete()
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
