package v1

import (
	"spiritFruit/app/models"
	"spiritFruit/app/models/projects"
	"spiritFruit/app/requests"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ProjectsController struct {
	BaseADMINController
}

// Index 短剧项目列表
// @Summary 短剧项目列表
// @Description 获取短剧项目列表，支持多种搜索条件
// @Tags Projects
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param per_page query int false "每页数量"
// @Param adminId query string false "归属用户ID(默认1)"
// @Param serialNo query string false "业务流水号"
// @Param title query string false "项目名称/短剧标题"
// @Param status query string false "状态"
// @Success 200 {object} response.Response{data=[]projects.Projects} "短剧项目列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/projects [get]
func (ctrl *ProjectsController) Index(c *gin.Context) {
	// 构建搜索条件
	where := ctrl.buildSearchConditions(c)

	// 获取分页参数
	perPage := 10
	if perPageStr := c.Query("pageSize"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	data, pager := projects.Paginate(c, perPage, where)
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
func (ctrl *ProjectsController) buildSearchConditions(c *gin.Context) map[string]interface{} {
	where := map[string]interface{}{}

	// 归属用户ID(默认1)搜索

	if adminId := strings.TrimSpace(c.Query("adminId")); adminId != "" {
		where["admin_id"] = adminId
	}

	// 业务流水号搜索

	if serialNo := strings.TrimSpace(c.Query("serialNo")); serialNo != "" {
		where["serial_no"] = serialNo
	}

	// 项目名称/短剧标题搜索

	if title := strings.TrimSpace(c.Query("title")); title != "" {
		where["title"] = title
	}

	// 状态搜索

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		where["status"] = status
	}

	return where
}

// Show 短剧项目详情
// @Summary 短剧项目详情
// @Description 获取短剧项目详情
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Projects ID"
// @Success 200 {object} response.Response{data=projects.Projects} "短剧项目详情"
// @Failure 404 {object} response.ErrorResponse "短剧项目不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/projects/{id} [get]
func (ctrl *ProjectsController) Show(c *gin.Context) {
	projectsModel := projects.Get(c.Param("id"))
	if projectsModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}
	response.JSON(c, gin.H{
		"code":    0,
		"data":    projectsModel,
		"message": "success",
	})
}

// Store 创建短剧项目
// @Summary 创建短剧项目
// @Description 创建新的短剧项目
// @Tags Projects
// @Accept json
// @Produce json
// @Param request body requests.ProjectsRequest true "短剧项目信息"
// @Success 201 {object} response.Response{data=projects.Projects} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/projects [post]
func (ctrl *ProjectsController) Store(c *gin.Context) {
	request := requests.ProjectsRequest{}
	if ok := requests.Validate(c, &request, requests.ProjectsSave); !ok {
		return
	}
	adminId := uint64(1)
	projectsModel := projects.Projects{
		AdminId:       &adminId,
		SerialNo:      &request.SerialNo,
		Title:         &request.Title,
		Description:   &request.Description,
		Style:         &request.Style,
		Status:        &request.Status,
		Image:         &request.Image,
		TotalDuration: &request.TotalDuration,
		Settings:      &request.Settings,
	}

	projectsModel.Create()
	if projectsModel.ID > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    projectsModel,
			"message": "success",
		})
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

// Update 更新短剧项目
// @Summary 更新短剧项目
// @Description 更新短剧项目信息
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Projects ID"
// @Param request body requests.ProjectsRequest true "短剧项目信息"
// @Success 200 {object} response.Response{data=projects.Projects} "更新成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "短剧项目不存在"
// @Failure 422 {object} response.ValidationErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/projects/{id} [put]
func (ctrl *ProjectsController) Update(c *gin.Context) {
	// 验证数据是否存在
	id := c.Param("id")
	existingProjects := projects.Get(id)
	if existingProjects.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	request := requests.ProjectsRequest{}
	if bindOk := requests.Validate(c, &request, requests.ProjectsSave); !bindOk {
		return
	}
	// 使用新的模型实例进行更新，避免关联对象的影响
	updateProjects := &projects.Projects{
		BaseModel: models.BaseModel{ID: existingProjects.ID},
	}

	// 赋值字段
	adminId := uint64(1)
	updateProjects.AdminId = &adminId
	updateProjects.SerialNo = &request.SerialNo
	updateProjects.Title = &request.Title
	updateProjects.Description = &request.Description
	updateProjects.Style = &request.Style
	updateProjects.Status = &request.Status
	updateProjects.Image = &request.Image
	updateProjects.TotalDuration = &request.TotalDuration
	updateProjects.Settings = &request.Settings
	updateProjects.UpdatedAt = time.Now()
	updateProjects.CreatedAt = existingProjects.CreatedAt

	// 执行更新
	result := database.DB.Save(updateProjects)

	if result.Error != nil {
		response.Abort500(c, "更新失败："+result.Error.Error())
		return
	}

	if result.RowsAffected > 0 {
		// 重新获取更新后的完整数据（包括关联）
		updatedProjects := projects.Get(id)
		response.JSON(c, gin.H{
			"code":    0,
			"data":    updatedProjects,
			"message": "success",
		})
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// Delete 删除短剧项目
// @Summary 删除短剧项目
// @Description 删除短剧项目
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Projects ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "短剧项目不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/projects/{id} [delete]
func (ctrl *ProjectsController) Delete(c *gin.Context) {
	projectsModel := projects.Get(c.Param("id"))
	if projectsModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	rowsAffected := projectsModel.Delete()
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
