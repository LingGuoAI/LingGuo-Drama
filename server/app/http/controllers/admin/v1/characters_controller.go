package v1

import (
    "spiritFruit/app/models/characters"
    "spiritFruit/app/models/projects"
    "spiritFruit/app/requests"
    "spiritFruit/pkg/response"
    "strconv"
    "strings"
    "time"
    "spiritFruit/app/models"
    "spiritFruit/pkg/database"
    "github.com/gin-gonic/gin"
)

type CharactersController struct {
    BaseADMINController
}

// Index 角色列表
// @Summary 角色列表
// @Description 获取角色列表，支持多种搜索条件
// @Tags Characters
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param per_page query int false "每页数量"
// @Param projectId query string false "所属项目ID"
// @Param name query string false "角色名"
// @Param roleType query string false "角色类型: main/supporting/minor"
// @Param gender query string false "性别(需从appearance解析或留空)"
// @Success 200 {object} response.Response{data=[]characters.Characters} "角色列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/characters [get]
func (ctrl *CharactersController) Index(c *gin.Context) {
    // 构建搜索条件
    where := ctrl.buildSearchConditions(c)

    // 获取分页参数
    perPage := 10
    if perPageStr := c.Query("per_page"); perPageStr != "" {
        if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
            perPage = pp
        }
    }

    data, pager := characters.Paginate(c, perPage, where)
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
func (ctrl *CharactersController) buildSearchConditions(c *gin.Context) map[string]interface{} {
    where := map[string]interface{}{}


    // 所属项目ID搜索
    
    if projectId := strings.TrimSpace(c.Query("projectId")); projectId != "" {
        where["project_id"] = projectId
    }
    

    // 角色名搜索
    
    if name := strings.TrimSpace(c.Query("name")); name != "" {
        where["name"] = name
    }
    

    // 角色类型: main/supporting/minor搜索
    
    if roleType := strings.TrimSpace(c.Query("roleType")); roleType != "" {
        where["role_type"] = roleType
    }
    

    // 性别(需从appearance解析或留空)搜索
    
    if gender := strings.TrimSpace(c.Query("gender")); gender != "" {
        where["gender"] = gender
    }
    


    return where
}
// GetProjectsSelectList 获取短剧项目选择列表
// @Summary 获取短剧项目选择列表
// @Description 获取短剧项目的简化列表，用于角色中的短剧项目选择
// @Tags Characters
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]projects.Projects} "短剧项目选择列表"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/characters/getProjectsSelectList [get]
func (ctrl *CharactersController) GetProjectsSelectList(c *gin.Context) {
    list:=projects.All()
    response.JSON(c, gin.H{
        "code":    0,
        "data":    list,
        "message": "success",
    })
}

// Show 角色详情
// @Summary 角色详情
// @Description 获取角色详情
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path string true "Characters ID"
// @Success 200 {object} response.Response{data=characters.Characters} "角色详情"
// @Failure 404 {object} response.ErrorResponse "角色不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/characters/{id} [get]
func (ctrl *CharactersController) Show(c *gin.Context) {
    charactersModel := characters.Get(c.Param("id"))
    if charactersModel.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }
    response.JSON(c, gin.H{
        "code":    0,
        "data":    charactersModel,
        "message": "success",
    })
}

// Store 创建角色
// @Summary 创建角色
// @Description 创建新的角色
// @Tags Characters
// @Accept json
// @Produce json
// @Param request body requests.CharactersRequest true "角色信息"
// @Success 201 {object} response.Response{data=characters.Characters} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/characters [post]
func (ctrl *CharactersController) Store(c *gin.Context) {
    request := requests.CharactersRequest{}
    if ok := requests.Validate(c, &request, requests.CharactersSave); !ok {
        return
    }
    charactersModel := characters.Characters{
        ProjectId: &request.ProjectId,
        Name: &request.Name,
        RoleType: &request.RoleType,
        Gender: &request.Gender,
        AgeGroup: &request.AgeGroup,
        Personality: &request.Personality,
        AppearanceDesc: &request.AppearanceDesc,
        VisualPrompt: &request.VisualPrompt,
        AvatarUrl: &request.AvatarUrl,
        VoiceId: &request.VoiceId,
    }

    charactersModel.Create()
    if charactersModel.ID > 0 {
        response.JSON(c, gin.H{
            "code":    0,
            "data":    charactersModel,
            "message": "success",
        })
    } else {
    response.Abort500(c, "创建失败，请稍后尝试~")
    }
}

// Update 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path string true "Characters ID"
// @Param request body requests.CharactersRequest true "角色信息"
// @Success 200 {object} response.Response{data=characters.Characters} "更新成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "角色不存在"
// @Failure 422 {object} response.ValidationErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/characters/{id} [put]
func (ctrl *CharactersController) Update(c *gin.Context) {
    // 验证数据是否存在
    id := c.Param("id")
    existingCharacters := characters.Get(id)
    if existingCharacters.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }

    request := requests.CharactersRequest{}
    if bindOk := requests.Validate(c, &request, requests.CharactersSave); !bindOk {
        return
    }
    // 使用新的模型实例进行更新，避免关联对象的影响
    updateCharacters := &characters.Characters{
        BaseModel: models.BaseModel{ID: existingCharacters.ID},
    }

    // 赋值字段
    updateCharacters.ProjectId = &request.ProjectId
    updateCharacters.Name = &request.Name
    updateCharacters.RoleType = &request.RoleType
    updateCharacters.Gender = &request.Gender
    updateCharacters.AgeGroup = &request.AgeGroup
    updateCharacters.Personality = &request.Personality
    updateCharacters.AppearanceDesc = &request.AppearanceDesc
    updateCharacters.VisualPrompt = &request.VisualPrompt
    updateCharacters.AvatarUrl = &request.AvatarUrl
    updateCharacters.VoiceId = &request.VoiceId
    updateCharacters.UpdatedAt = time.Now()
    updateCharacters.CreatedAt = existingCharacters.CreatedAt

    // 执行更新
    result := database.DB.Save(updateCharacters)

    if result.Error != nil {
        response.Abort500(c, "更新失败：" + result.Error.Error())
        return
    }

    if result.RowsAffected > 0 {
        // 重新获取更新后的完整数据（包括关联）
        updatedCharacters := characters.Get(id)
        response.JSON(c, gin.H{
            "code":    0,
            "data":    updatedCharacters,
            "message": "success",
        })
    } else {
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

// Delete 删除角色
// @Summary 删除角色
// @Description 删除角色
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path string true "Characters ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "角色不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/characters/{id} [delete]
func (ctrl *CharactersController) Delete(c *gin.Context) {
    charactersModel := characters.Get(c.Param("id"))
    if charactersModel.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }

    rowsAffected := charactersModel.Delete()
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