package sys_base_menus

import (
    "spiritFruit/pkg/app"
    "spiritFruit/pkg/database"
    "spiritFruit/pkg/paginator"
    "github.com/gin-gonic/gin"
)


func Get(idStr string) (sysBaseMenus SysBaseMenus) {
    database.DB.Where("id", idStr).First(&sysBaseMenus)
    return
}

func GetBy(field string, value uint64) (sysBaseMenus SysBaseMenus) {
    database.DB.Where(map[string]interface{}{field: value}).First(&sysBaseMenus)
    return
}

func GetByWhereMap(where map[string]interface{}) (sysBaseMenus SysBaseMenus) {
	database.DB.Where(where).First(&sysBaseMenus)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (sysBaseMenus []SysBaseMenus) {
	database.DB.Where(where).Find(&sysBaseMenus)
	return
}

func All() (sysBaseMenuses []SysBaseMenus) {
    database.DB.Find(&sysBaseMenuses)
    return
}

func IsExist(field string, value uint64) bool {
    var count int64
    database.DB.Model(&SysBaseMenus{}).Where(map[string]interface{}{field: value}).Count(&count)
    return count > 0
}

// Paginate 系统菜单分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (sysBaseMenuses []SysBaseMenus, paging paginator.Paging) {
    paging = paginator.Paginate(
        c,
        database.DB.Model(SysBaseMenus{}),
        &sysBaseMenuses,
        app.V1URL(database.TableName(&SysBaseMenus{})),
        perPage,
        filters,
    )
    return
}