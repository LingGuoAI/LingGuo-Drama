package v1

import (
    "spiritFruit/app/models/sys_base_menus"
    "spiritFruit/app/requests"
    "spiritFruit/pkg/response"
    "strconv"
    "strings"
    "time"
    "spiritFruit/app/models"
    "spiritFruit/pkg/database"
    "github.com/gin-gonic/gin"
)

type SysBaseMenusesController struct {
    BaseADMINController
}

// Index 系统菜单列表
// @Summary 系统菜单列表
// @Description 获取系统菜单列表，支持多种搜索条件
// @Tags SysBaseMenuses
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param per_page query int false "每页数量"
// @Param parentId query string false "父菜单ID"
// @Param name query string false "路由名称"
// @Param title query string false "菜单标题"
// @Success 200 {object} response.Response{data=[]sys_base_menus.SysBaseMenus} "系统菜单列表"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sys_base_menuses [get]
func (ctrl *SysBaseMenusesController) Index(c *gin.Context) {
    // 构建搜索条件
    where := ctrl.buildSearchConditions(c)

    // 获取分页参数
    perPage := 10
    if perPageStr := c.Query("per_page"); perPageStr != "" {
        if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
            perPage = pp
        }
    }

    data, pager := sys_base_menus.Paginate(c, perPage, where)
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
func (ctrl *SysBaseMenusesController) buildSearchConditions(c *gin.Context) map[string]interface{} {
    where := map[string]interface{}{}


    // 父菜单ID搜索
    
    if parentId := strings.TrimSpace(c.Query("parentId")); parentId != "" {
        where["parent_id"] = parentId
    }
    

    // 路由名称搜索
    
    if name := strings.TrimSpace(c.Query("name")); name != "" {
        where["name"] = name
    }
    

    // 菜单标题搜索
    
    if title := strings.TrimSpace(c.Query("title")); title != "" {
        where["title"] = title
    }
    


    return where
}
// GetSysBaseMenusTreeList 获取系统菜单树
// @Summary 获取系统菜单树
// @Description 获取完整的系统菜单树形结构，用于树形选择器和树形表格显示
// @Tags SysBaseMenuses
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]sys_base_menus.SysBaseMenus} "系统菜单树"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sys_base_menuses/getSysBaseMenusTreeList [get]
func (ctrl *SysBaseMenusesController) GetSysBaseMenusTreeList(c *gin.Context) {
    
    tree := sys_base_menus.GetSysBaseMenusTree()
        response.JSON(c, gin.H{
        "code":    0,
        "data":    tree,
        "message": "success",
    })
}
// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取系统菜单树形结构列表
// @Tags SysBaseMenuses
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=sys_base_menus.MenuListResult} "菜单列表"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sys_base_menuses/getMenuList [get]
func (ctrl *SysBaseMenusesController) GetMenuList(c *gin.Context) {
    menuTree, err := sys_base_menus.BuildMenuTree()
    if err != nil {
        response.JSON(c, gin.H{
            "code":    500,
            "message": "获取菜单列表失败",
            "data":    nil,
        })
        return
    }

    result := sys_base_menus.MenuListResult{
        List: menuTree,
    }

    response.JSON(c, gin.H{
        "code":    0,
        "message": "获取菜单列表成功",
        "data":    result,
    })
}

// GetAuthorityId 按id结构获取角色
type GetAuthorityId struct {
    AuthorityId uint64 `json:"authorityId" form:"authorityId"` // 角色ID
}

// Show 系统菜单详情
// @Summary 系统菜单详情
// @Description 获取系统菜单详情
// @Tags SysBaseMenuses
// @Accept json
// @Produce json
// @Param id path string true "SysBaseMenus ID"
// @Success 200 {object} response.Response{data=sys_base_menus.SysBaseMenus} "系统菜单详情"
// @Failure 404 {object} response.ErrorResponse "系统菜单不存在"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sys_base_menuses/{id} [get]
func (ctrl *SysBaseMenusesController) Show(c *gin.Context) {
    sysBaseMenusModel := sys_base_menus.Get(c.Param("id"))
    if sysBaseMenusModel.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }
    response.JSON(c, gin.H{
        "code":    0,
        "data":    sysBaseMenusModel,
        "message": "success",
    })
}

// Store 创建系统菜单
// @Summary 创建系统菜单
// @Description 创建新的系统菜单
// @Tags SysBaseMenuses
// @Accept json
// @Produce json
// @Param request body requests.SysBaseMenusRequest true "系统菜单信息"
// @Success 201 {object} response.Response{data=sys_base_menus.SysBaseMenus} "创建成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 422 {object} response.ErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sys_base_menuses [post]
func (ctrl *SysBaseMenusesController) Store(c *gin.Context) {
    request := requests.SysBaseMenusRequest{}
    if ok := requests.Validate(c, &request, requests.SysBaseMenusSave); !ok {
        return
    }
    sysBaseMenusModel := sys_base_menus.SysBaseMenus{
        ParentId: &request.ParentId,
        Path: &request.Path,
        Name: &request.Name,
        Hidden: &request.Hidden,
        Component: &request.Component,
        Sort: &request.Sort,
        Title: &request.Title,
        Icon: &request.Icon,
    }

    sysBaseMenusModel.Create()
    if sysBaseMenusModel.ID > 0 {
        response.JSON(c, gin.H{
            "code":    0,
            "data":    sysBaseMenusModel,
            "message": "success",
        })
    } else {
    response.Abort500(c, "创建失败，请稍后尝试~")
    }
}

// Update 更新系统菜单
// @Summary 更新系统菜单
// @Description 更新系统菜单信息
// @Tags SysBaseMenuses
// @Accept json
// @Produce json
// @Param id path string true "SysBaseMenus ID"
// @Param request body requests.SysBaseMenusRequest true "系统菜单信息"
// @Success 200 {object} response.Response{data=sys_base_menus.SysBaseMenus} "更新成功"
// @Failure 400 {object} response.ErrorResponse "参数错误"
// @Failure 404 {object} response.ErrorResponse "系统菜单不存在"
// @Failure 422 {object} response.ValidationErrorResponse "验证失败"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sys_base_menuses/{id} [put]
func (ctrl *SysBaseMenusesController) Update(c *gin.Context) {
    // 验证数据是否存在
    id := c.Param("id")
    existingSysBaseMenus := sys_base_menus.Get(id)
    if existingSysBaseMenus.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }

    request := requests.SysBaseMenusRequest{}
    if bindOk := requests.Validate(c, &request, requests.SysBaseMenusSave); !bindOk {
        return
    }
    // 使用新的模型实例进行更新，避免关联对象的影响
    updateSysBaseMenus := &sys_base_menus.SysBaseMenus{
        BaseModel: models.BaseModel{ID: existingSysBaseMenus.ID},
    }

    // 赋值字段
    updateSysBaseMenus.ParentId = &request.ParentId
    updateSysBaseMenus.Path = &request.Path
    updateSysBaseMenus.Name = &request.Name
    updateSysBaseMenus.Hidden = &request.Hidden
    updateSysBaseMenus.Component = &request.Component
    updateSysBaseMenus.Sort = &request.Sort
    updateSysBaseMenus.Title = &request.Title
    updateSysBaseMenus.Icon = &request.Icon
    updateSysBaseMenus.UpdatedAt = time.Now()
    updateSysBaseMenus.CreatedAt = existingSysBaseMenus.CreatedAt

    // 执行更新
    result := database.DB.Save(updateSysBaseMenus)

    if result.Error != nil {
        response.Abort500(c, "更新失败：" + result.Error.Error())
        return
    }

    if result.RowsAffected > 0 {
        // 重新获取更新后的完整数据（包括关联）
        updatedSysBaseMenus := sys_base_menus.Get(id)
        response.JSON(c, gin.H{
            "code":    0,
            "data":    updatedSysBaseMenus,
            "message": "success",
        })
    } else {
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

// Delete 删除系统菜单
// @Summary 删除系统菜单
// @Description 删除系统菜单
// @Tags SysBaseMenuses
// @Accept json
// @Produce json
// @Param id path string true "SysBaseMenus ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 404 {object} response.ErrorResponse "系统菜单不存在"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/v1/sys_base_menuses/{id} [delete]
func (ctrl *SysBaseMenusesController) Delete(c *gin.Context) {
    sysBaseMenusModel := sys_base_menus.Get(c.Param("id"))
    if sysBaseMenusModel.ID == 0 {
        response.JSON(c, gin.H{
            "code":    404,
            "message": "数据不存在",
            "data":    nil,
        })
        return
    }

    rowsAffected := sysBaseMenusModel.Delete()
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