package v1

import (
    "spiritFruit/app/models/admins"
    "spiritFruit/app/requests"
    "spiritFruit/pkg/response"
    "strconv"
    "strings"
    "time"
    "spiritFruit/app/models"
    "spiritFruit/pkg/database"
    "github.com/gin-gonic/gin"
)

type AdminsController struct {
    BaseADMINController
}

// Index 系统管理员列表
// @Summary 系统管理员列表
// @Description 获取系统管理员列表，支持多种搜索条件
// @Tags Admins
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param per_page query int false "每页数量"
// @Param username query string false "用户名"
// @Param mobile query string false "手机号"
// @Param email query string false "邮箱"
// @Success 200 {object} response.Response{data=[]admins.Admins} "系统管理员列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/admins [get]
func (ctrl *AdminsController) Index(c *gin.Context) {
    // 构建搜索条件
    where := ctrl.buildSearchConditions(c)

    // 获取分页参数
    perPage := 10
    if perPageStr := c.Query("per_page"); perPageStr != "" {
        if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
            perPage = pp
        }
    }

    data, pager := admins.Paginate(c, perPage, where)
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
func (ctrl *AdminsController) buildSearchConditions(c *gin.Context) map[string]interface{} {
    where := map[string]interface{}{}


    // 用户名搜索
    
    if username := strings.TrimSpace(c.Query("username")); username != "" {
        where["username LIKE ?"] = "%" + username + "%"
    }
    

    // 手机号搜索
    
    if mobile := strings.TrimSpace(c.Query("mobile")); mobile != "" {
        where["mobile LIKE ?"] = "%" + mobile + "%"
    }
    

    // 邮箱搜索
    
    if email := strings.TrimSpace(c.Query("email")); email != "" {
        where["email LIKE ?"] = "%" + email + "%"
    }
    


    return where
}

// Show 系统管理员详情
// @Summary 系统管理员详情
// @Description 获取系统管理员详情
// @Tags Admins
// @Accept json
// @Produce json
// @Param id path string true "Admins ID"
// @Success 200 {object} response.Response{data=admins.Admins} "系统管理员详情"
// @Failure 404 {object} response.ErrorResponse "系统管理员不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/admins/{id} [get]
func (ctrl *AdminsController) Show(c *gin.Context) {
    adminsModel := admins.Get(c.Param("id"))
    if adminsModel.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }
    response.JSON(c, gin.H{
        "code":    0,
        "data":    adminsModel,
        "message": "success",
    })
}

// Store 创建系统管理员
// @Summary 创建系统管理员
// @Description 创建新的系统管理员
// @Tags Admins
// @Accept json
// @Produce json
// @Param request body requests.AdminsRequest true "系统管理员信息"
// @Success 201 {object} response.Response{data=admins.Admins} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/admins [post]
func (ctrl *AdminsController) Store(c *gin.Context) {
    request := requests.AdminsRequest{}
    if ok := requests.Validate(c, &request, requests.AdminsSave); !ok {
        return
    }
    adminsModel := admins.Admins{
        Username: &request.Username,
        Mobile: &request.Mobile,
        Password: &request.Password,
        Email: &request.Email,
        AuthorityId: &request.AuthorityId,
    }

    adminsModel.Create()
    if adminsModel.ID > 0 {
        response.JSON(c, gin.H{
            "code":    0,
            "data":    adminsModel,
            "message": "success",
        })
    } else {
    response.Abort500(c, "创建失败，请稍后尝试~")
    }
}

// Update 更新系统管理员
// @Summary 更新系统管理员
// @Description 更新系统管理员信息
// @Tags Admins
// @Accept json
// @Produce json
// @Param id path string true "Admins ID"
// @Param request body requests.AdminsRequest true "系统管理员信息"
// @Success 200 {object} response.Response{data=admins.Admins} "更新成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "系统管理员不存在"
// @Failure 422 {object} response.ValidationErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/admins/{id} [put]
func (ctrl *AdminsController) Update(c *gin.Context) {
    // 验证数据是否存在
    id := c.Param("id")
    existingAdmins := admins.Get(id)
    if existingAdmins.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }

    request := requests.AdminsRequest{}
    if bindOk := requests.Validate(c, &request, requests.AdminsSave); !bindOk {
        return
    }
    // 使用新的模型实例进行更新，避免关联对象的影响
    updateAdmins := &admins.Admins{
        BaseModel: models.BaseModel{ID: existingAdmins.ID},
    }

    // 赋值字段
    updateAdmins.Username = &request.Username
    updateAdmins.Mobile = &request.Mobile
    updateAdmins.Password = &request.Password
    updateAdmins.Email = &request.Email
    updateAdmins.AuthorityId = &request.AuthorityId
    updateAdmins.UpdatedAt = time.Now()
    updateAdmins.CreatedAt = existingAdmins.CreatedAt

    // 执行更新
    result := database.DB.Save(updateAdmins)

    if result.Error != nil {
        response.Abort500(c, "更新失败：" + result.Error.Error())
        return
    }

    if result.RowsAffected > 0 {
        // 重新获取更新后的完整数据（包括关联）
        updatedAdmins := admins.Get(id)
        response.JSON(c, gin.H{
            "code":    0,
            "data":    updatedAdmins,
            "message": "success",
        })
    } else {
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

// Delete 删除系统管理员
// @Summary 删除系统管理员
// @Description 删除系统管理员
// @Tags Admins
// @Accept json
// @Produce json
// @Param id path string true "Admins ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "系统管理员不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/admins/{id} [delete]
func (ctrl *AdminsController) Delete(c *gin.Context) {
    adminsModel := admins.Get(c.Param("id"))
    if adminsModel.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }

    rowsAffected := adminsModel.Delete()
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