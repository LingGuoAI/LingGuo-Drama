package v1

import (
    "spiritFruit/app/models/shots"
    "spiritFruit/app/models/projects"
    "spiritFruit/app/models/scripts"
    "spiritFruit/app/requests"
    "spiritFruit/pkg/response"
    "strconv"
    "strings"
    "time"
    "spiritFruit/app/models"
    "spiritFruit/pkg/database"
    "github.com/gin-gonic/gin"
)

type ShotsController struct {
    BaseADMINController
}

// Index 镜头表列表
// @Summary 镜头表列表
// @Description 获取镜头表列表，支持多种搜索条件
// @Tags Shots
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param per_page query int false "每页数量"
// @Param projectId query string false "所属项目ID"
// @Param scriptId query string false "所属剧本/分集ID"
// @Param sequenceNo query string false "镜头序号"
// @Success 200 {object} response.Response{data=[]shots.Shots} "镜头表列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shots [get]
func (ctrl *ShotsController) Index(c *gin.Context) {
    // 构建搜索条件
    where := ctrl.buildSearchConditions(c)

    // 获取分页参数
    perPage := 10
    if perPageStr := c.Query("per_page"); perPageStr != "" {
        if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
            perPage = pp
        }
    }

    data, pager := shots.Paginate(c, perPage, where)
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
func (ctrl *ShotsController) buildSearchConditions(c *gin.Context) map[string]interface{} {
    where := map[string]interface{}{}


    // 所属项目ID搜索
    
    if projectId := strings.TrimSpace(c.Query("projectId")); projectId != "" {
        where["project_id"] = projectId
    }
    

    // 所属剧本/分集ID搜索
    
    if scriptId := strings.TrimSpace(c.Query("scriptId")); scriptId != "" {
        where["script_id"] = scriptId
    }
    

    // 镜头序号搜索
    
    if sequenceNo := strings.TrimSpace(c.Query("sequenceNo")); sequenceNo != "" {
        where["sequence_no"] = sequenceNo
    }
    


    return where
}
// GetProjectsSelectList 获取短剧项目选择列表
// @Summary 获取短剧项目选择列表
// @Description 获取短剧项目的简化列表，用于镜头表中的短剧项目选择
// @Tags Shots
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]projects.Projects} "短剧项目选择列表"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shots/getProjectsSelectList [get]
func (ctrl *ShotsController) GetProjectsSelectList(c *gin.Context) {
    list:=projects.All()
    response.JSON(c, gin.H{
        "code":    0,
        "data":    list,
        "message": "success",
    })
}
// GetScriptsSelectList 获取剧本选择列表
// @Summary 获取剧本选择列表
// @Description 获取剧本的简化列表，用于镜头表中的剧本选择
// @Tags Shots
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]scripts.Scripts} "剧本选择列表"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shots/getScriptsSelectList [get]
func (ctrl *ShotsController) GetScriptsSelectList(c *gin.Context) {
    list:=scripts.All()
    response.JSON(c, gin.H{
        "code":    0,
        "data":    list,
        "message": "success",
    })
}

// Show 镜头表详情
// @Summary 镜头表详情
// @Description 获取镜头表详情
// @Tags Shots
// @Accept json
// @Produce json
// @Param id path string true "Shots ID"
// @Success 200 {object} response.Response{data=shots.Shots} "镜头表详情"
// @Failure 404 {object} response.ErrorResponse "镜头表不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shots/{id} [get]
func (ctrl *ShotsController) Show(c *gin.Context) {
    shotsModel := shots.Get(c.Param("id"))
    if shotsModel.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }
    response.JSON(c, gin.H{
        "code":    0,
        "data":    shotsModel,
        "message": "success",
    })
}

// Store 创建镜头表
// @Summary 创建镜头表
// @Description 创建新的镜头表
// @Tags Shots
// @Accept json
// @Produce json
// @Param request body requests.ShotsRequest true "镜头表信息"
// @Success 201 {object} response.Response{data=shots.Shots} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shots [post]
func (ctrl *ShotsController) Store(c *gin.Context) {
    request := requests.ShotsRequest{}
    if ok := requests.Validate(c, &request, requests.ShotsSave); !ok {
        return
    }
    shotsModel := shots.Shots{
        ProjectId: &request.ProjectId,
        ScriptId: &request.ScriptId,
        SequenceNo: &request.SequenceNo,
        ShotType: &request.ShotType,
        CameraMovement: &request.CameraMovement,
        Angle: &request.Angle,
        Dialogue: &request.Dialogue,
        VisualDesc: &request.VisualDesc,
        Atmosphere: &request.Atmosphere,
        ImagePrompt: &request.ImagePrompt,
        VideoPrompt: &request.VideoPrompt,
        AudioPrompt: &request.AudioPrompt,
        ImageUrl: &request.ImageUrl,
        VideoUrl: &request.VideoUrl,
        AudioUrl: &request.AudioUrl,
        DurationMs: &request.DurationMs,
        Status: &request.Status,
    }

    shotsModel.Create()
    if shotsModel.ID > 0 {
        response.JSON(c, gin.H{
            "code":    0,
            "data":    shotsModel,
            "message": "success",
        })
    } else {
    response.Abort500(c, "创建失败，请稍后尝试~")
    }
}

// Update 更新镜头表
// @Summary 更新镜头表
// @Description 更新镜头表信息
// @Tags Shots
// @Accept json
// @Produce json
// @Param id path string true "Shots ID"
// @Param request body requests.ShotsRequest true "镜头表信息"
// @Success 200 {object} response.Response{data=shots.Shots} "更新成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "镜头表不存在"
// @Failure 422 {object} response.ValidationErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shots/{id} [put]
func (ctrl *ShotsController) Update(c *gin.Context) {
    // 验证数据是否存在
    id := c.Param("id")
    existingShots := shots.Get(id)
    if existingShots.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }

    request := requests.ShotsRequest{}
    if bindOk := requests.Validate(c, &request, requests.ShotsSave); !bindOk {
        return
    }
    // 使用新的模型实例进行更新，避免关联对象的影响
    updateShots := &shots.Shots{
        BaseModel: models.BaseModel{ID: existingShots.ID},
    }

    // 赋值字段
    updateShots.ProjectId = &request.ProjectId
    updateShots.ScriptId = &request.ScriptId
    updateShots.SequenceNo = &request.SequenceNo
    updateShots.ShotType = &request.ShotType
    updateShots.CameraMovement = &request.CameraMovement
    updateShots.Angle = &request.Angle
    updateShots.Dialogue = &request.Dialogue
    updateShots.VisualDesc = &request.VisualDesc
    updateShots.Atmosphere = &request.Atmosphere
    updateShots.ImagePrompt = &request.ImagePrompt
    updateShots.VideoPrompt = &request.VideoPrompt
    updateShots.AudioPrompt = &request.AudioPrompt
    updateShots.ImageUrl = &request.ImageUrl
    updateShots.VideoUrl = &request.VideoUrl
    updateShots.AudioUrl = &request.AudioUrl
    updateShots.DurationMs = &request.DurationMs
    updateShots.Status = &request.Status
    updateShots.UpdatedAt = time.Now()
    updateShots.CreatedAt = existingShots.CreatedAt

    // 执行更新
    result := database.DB.Save(updateShots)

    if result.Error != nil {
        response.Abort500(c, "更新失败：" + result.Error.Error())
        return
    }

    if result.RowsAffected > 0 {
        // 重新获取更新后的完整数据（包括关联）
        updatedShots := shots.Get(id)
        response.JSON(c, gin.H{
            "code":    0,
            "data":    updatedShots,
            "message": "success",
        })
    } else {
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

// Delete 删除镜头表
// @Summary 删除镜头表
// @Description 删除镜头表
// @Tags Shots
// @Accept json
// @Produce json
// @Param id path string true "Shots ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "镜头表不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/shots/{id} [delete]
func (ctrl *ShotsController) Delete(c *gin.Context) {
    shotsModel := shots.Get(c.Param("id"))
    if shotsModel.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }

    rowsAffected := shotsModel.Delete()
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