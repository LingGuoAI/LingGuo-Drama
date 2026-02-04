package characters

import (
    "gorm.io/gorm"
    "strings"
    "spiritFruit/pkg/app"
    "spiritFruit/pkg/database"
    "spiritFruit/pkg/paginator"
    "github.com/gin-gonic/gin"
)


func Get(idStr string) (characters Characters) {
    database.DB.Where("id", idStr).
    Preload("Projectss", func(db *gorm.DB) *gorm.DB {
    // 只加载关联表的必要字段
    fields := []string{"id","admin_id","serial_no","title","description","status","image","total_duration" }
        return db.Select(strings.Join(fields, ", "))
    }).First(&characters)
    return
}

func GetBy(field string, value uint64) (characters Characters) {
    database.DB.Where(map[string]interface{}{field: value}).First(&characters)
    return
}

func GetByWhereMap(where map[string]interface{}) (characters Characters) {
	database.DB.Where(where).First(&characters)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (characters []Characters) {
	database.DB.Where(where).Find(&characters)
	return
}

func All() (characters []Characters) {
    database.DB.Find(&characters)
    return
}

func IsExist(field string, value uint64) bool {
    var count int64
    database.DB.Model(&Characters{}).Where(map[string]interface{}{field: value}).Count(&count)
    return count > 0
}

// Paginate 角色分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (characters []Characters, paging paginator.Paging) {
    // 构建带关联预加载的查询
    query := database.DB.Model(Characters{}).
        Preload("Projectss", func(db *gorm.DB) *gorm.DB {
            // 只加载关联表的必要字段
            fields := []string{"id","admin_id","serial_no","title","description","status","image","total_duration" }
            return db.Select(strings.Join(fields, ", "))
    })

    // 应用过滤条件
    for key, value := range filters {
        if key == "ORDER" {
            query = query.Order(value.(string))
        } else if key == "LIMIT" {
            query = query.Limit(value.(int))
        } else {
            query = query.Where(key, value)
        }
    }

    // 使用自定义查询进行分页
    paging = paginator.PaginateCustomQuery(
        c,
        query,
        &characters,
        app.V1URL(database.TableName(&Characters{})),
        perPage,
        )
    return
}