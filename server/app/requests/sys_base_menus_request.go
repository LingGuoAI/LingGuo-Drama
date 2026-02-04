package requests

import (
    "github.com/gin-gonic/gin"
    "github.com/thedevsaddam/govalidator"
)

// SysBaseMenusRequest 系统菜单请求结构体
type SysBaseMenusRequest struct {
    ParentId uint64 `valid:"parentId" json:"parentId"` // 父菜单ID
    Path string `valid:"path" json:"path" ` // 路由路径
    Name string `valid:"name" json:"name" ` // 路由名称
    Hidden int8 `valid:"hidden" json:"hidden"` // 是否隐藏
    Component string `valid:"component" json:"component"` // 组件路径
    Sort uint64 `valid:"sort" json:"sort"` // 排序
    Title string `valid:"title" json:"title" ` // 菜单标题
    Icon string `valid:"icon" json:"icon"` // 菜单图标
}

// SysBaseMenusSave 系统菜单保存时的验证规则
func SysBaseMenusSave(data interface{}, c *gin.Context) map[string][]string {
    rules := govalidator.MapData{
        "parentId": []string{"numeric"},
        "path": []string{"required", "max:255", "min:1"},
        "name": []string{"required", "max:255", "min:1"},
        "hidden": []string{"numeric"},
        "component": []string{"max:255"},
        "sort": []string{"numeric"},
        "title": []string{"required", "max:255", "min:1"},
        "icon": []string{"max:255"},
    }

    messages := govalidator.MapData{
        "parentId": []string{
             "numeric:父菜单ID必须为数字",
        },
        "path": []string{
             "required:路由路径为必填项",
             "min:路由路径长度不能为空",
             "max:路由路径长度不能超过 255 个字符",
        },
        "name": []string{
             "required:路由名称为必填项",
             "min:路由名称长度不能为空",
             "max:路由名称长度不能超过 255 个字符",
        },
        "hidden": []string{
             "numeric:是否隐藏必须为数字",
        },
        "component": []string{
             "max:组件路径长度不能超过 255 个字符",
        },
        "sort": []string{
             "numeric:排序必须为数字",
        },
        "title": []string{
             "required:菜单标题为必填项",
             "min:菜单标题长度不能为空",
             "max:菜单标题长度不能超过 255 个字符",
        },
        "icon": []string{
             "max:菜单图标长度不能超过 255 个字符",
        },
    }

    return validate(data, rules, messages)
}