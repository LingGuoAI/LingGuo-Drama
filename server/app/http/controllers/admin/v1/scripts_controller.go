package v1

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/app/models"
	"spiritFruit/app/models/projects"
	"spiritFruit/app/models/scripts"
	"spiritFruit/app/requests"
	"spiritFruit/app/services"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
	"strconv"
	"strings"
	"time"
)

type ScriptsController struct {
	BaseADMINController
}

// Index 剧本列表
// @Summary 剧本列表
// @Description 获取剧本列表，支持多种搜索条件
// @Tags Scripts
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param per_page query int false "每页数量"
// @Param projectId query string false "所属项目ID"
// @Param title query string false "分集标题"
// @Param episodeNo query string false "第几集"
// @Param isLocked query string false "是否定稿"
// @Success 200 {object} response.Response{data=[]scripts.Scripts} "剧本列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scripts [get]
func (ctrl *ScriptsController) Index(c *gin.Context) {
	// 构建搜索条件
	where := ctrl.buildSearchConditions(c)

	// 获取分页参数
	perPage := 10
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	data, pager := scripts.Paginate(c, perPage, where)
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
func (ctrl *ScriptsController) buildSearchConditions(c *gin.Context) map[string]interface{} {
	where := map[string]interface{}{}

	// 所属项目ID搜索

	if projectId := strings.TrimSpace(c.Query("projectId")); projectId != "" {
		where["project_id"] = projectId
	}

	// 分集标题搜索

	if title := strings.TrimSpace(c.Query("title")); title != "" {
		where["title LIKE ?"] = "%" + title + "%"
	}

	// 第几集搜索

	if episodeNo := strings.TrimSpace(c.Query("episodeNo")); episodeNo != "" {
		where["episode_no"] = episodeNo
	}

	// 是否定稿搜索

	if isLocked := strings.TrimSpace(c.Query("isLocked")); isLocked != "" {
		where["is_locked"] = isLocked
	}

	return where
}

// GetProjectsSelectList 获取短剧项目选择列表
// @Summary 获取短剧项目选择列表
// @Description 获取短剧项目的简化列表，用于剧本中的短剧项目选择
// @Tags Scripts
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]projects.Projects} "短剧项目选择列表"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scripts/getProjectsSelectList [get]
func (ctrl *ScriptsController) GetProjectsSelectList(c *gin.Context) {
	list := projects.All()
	response.JSON(c, gin.H{
		"code":    0,
		"data":    list,
		"message": "success",
	})
}

// Generate 生成剧本 (异步任务)
func (ctrl *ScriptsController) Generate(c *gin.Context) {
	// 1. 定义请求参数结构
	type GenerateReq struct {
		ProjectID uint64 `json:"projectId" binding:"required"` // 项目ID
		ScriptID  uint64 `json:"scriptId" binding:"required"`  // 剧本ID (如果是新剧本，先创建空记录拿到ID传过来)
		Prompt    string `json:"prompt" binding:"required"`    // AI 提示词/大纲
	}

	var req GenerateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, "参数错误")
		return
	}

	// 2. 调用 TaskService 创建异步任务
	taskService := new(services.TaskService)
	task, err := taskService.CreateScriptGenerationTask(req.ProjectID, req.ScriptID, req.Prompt)

	if err != nil {
		response.Abort500(c, "任务启动失败: "+err.Error())
		return
	}

	// 3. 返回任务 ID 给前端进行轮询
	response.Data(c, gin.H{
		"taskId": task.ID,     // 前端拿这个 ID 去轮询 /tasks/:id
		"status": task.Status, // 初始状态 (0-Pending)
		"msg":    "剧本生成任务已提交，正在后台处理中...",
	})
}

// Show 剧本详情
// @Summary 剧本详情
// @Description 获取剧本详情
// @Tags Scripts
// @Accept json
// @Produce json
// @Param id path string true "Scripts ID"
// @Success 200 {object} response.Response{data=scripts.Scripts} "剧本详情"
// @Failure 404 {object} response.ErrorResponse "剧本不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scripts/{id} [get]
func (ctrl *ScriptsController) Show(c *gin.Context) {
	scriptsModel := scripts.Get(c.Param("id"))
	if scriptsModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}
	response.JSON(c, gin.H{
		"code":    0,
		"data":    scriptsModel,
		"message": "success",
	})
}

// Store 创建剧本
// @Summary 创建剧本
// @Description 创建新的剧本
// @Tags Scripts
// @Accept json
// @Produce json
// @Param request body requests.ScriptsRequest true "剧本信息"
// @Success 201 {object} response.Response{data=scripts.Scripts} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scripts [post]
func (ctrl *ScriptsController) Store(c *gin.Context) {
	request := requests.ScriptsRequest{}
	if ok := requests.Validate(c, &request, requests.ScriptsSave); !ok {
		return
	}
	scriptsModel := scripts.Scripts{
		ProjectId: &request.ProjectId,
		Title:     &request.Title,
		Content:   &request.Content,
		Outline:   &request.Outline,
		EpisodeNo: &request.EpisodeNo,
		IsLocked:  &request.IsLocked,
	}

	scriptsModel.Create()
	if scriptsModel.ID > 0 {
		response.JSON(c, gin.H{
			"code":    0,
			"data":    scriptsModel,
			"message": "success",
		})
	} else {
		response.Abort500(c, "创建失败，请稍后尝试~")
	}
}

// Update 更新剧本
// @Summary 更新剧本
// @Description 更新剧本信息
// @Tags Scripts
// @Accept json
// @Produce json
// @Param id path string true "Scripts ID"
// @Param request body requests.ScriptsRequest true "剧本信息"
// @Success 200 {object} response.Response{data=scripts.Scripts} "更新成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "剧本不存在"
// @Failure 422 {object} response.ValidationErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scripts/{id} [put]
func (ctrl *ScriptsController) Update(c *gin.Context) {
	// 验证数据是否存在
	id := c.Param("id")
	existingScripts := scripts.Get(id)
	if existingScripts.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	request := requests.ScriptsRequest{}
	if bindOk := requests.Validate(c, &request, requests.ScriptsSave); !bindOk {
		return
	}
	// 使用新的模型实例进行更新，避免关联对象的影响
	updateScripts := &scripts.Scripts{
		BaseModel: models.BaseModel{ID: existingScripts.ID},
	}

	// 赋值字段
	updateScripts.ProjectId = &request.ProjectId
	updateScripts.Title = &request.Title
	updateScripts.Content = &request.Content
	updateScripts.Outline = &request.Outline
	updateScripts.EpisodeNo = &request.EpisodeNo
	updateScripts.IsLocked = &request.IsLocked
	updateScripts.UpdatedAt = time.Now()
	updateScripts.CreatedAt = existingScripts.CreatedAt

	// 执行更新
	result := database.DB.Save(updateScripts)

	if result.Error != nil {
		response.Abort500(c, "更新失败："+result.Error.Error())
		return
	}

	if result.RowsAffected > 0 {
		// 重新获取更新后的完整数据（包括关联）
		updatedScripts := scripts.Get(id)
		response.JSON(c, gin.H{
			"code":    0,
			"data":    updatedScripts,
			"message": "success",
		})
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}

// Delete 删除剧本
// @Summary 删除剧本
// @Description 删除剧本
// @Tags Scripts
// @Accept json
// @Produce json
// @Param id path string true "Scripts ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "剧本不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/scripts/{id} [delete]
func (ctrl *ScriptsController) Delete(c *gin.Context) {
	scriptsModel := scripts.Get(c.Param("id"))
	if scriptsModel.ID == 0 {
		response.JSON(c, gin.H{
			"code":    404,
			"message": "数据不存在",
			"data":    nil,
		})
		return
	}

	rowsAffected := scriptsModel.Delete()
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
