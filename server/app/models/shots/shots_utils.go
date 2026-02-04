package shots

import (
    "gorm.io/gorm"
    "strings"
    "spiritFruit/pkg/app"
    "spiritFruit/pkg/database"
    "spiritFruit/pkg/paginator"
    "github.com/gin-gonic/gin"
)


func Get(idStr string) (shots Shots) {
    database.DB.Where("id", idStr).
    Preload("Projectss", func(db *gorm.DB) *gorm.DB {
    // 只加载关联表的必要字段
    fields := []string{"id","admin_id","serial_no","status","image" }
        return db.Select(strings.Join(fields, ", "))
    }).
    Preload("Scriptss", func(db *gorm.DB) *gorm.DB {
    // 只加载关联表的必要字段
    fields := []string{"id","project_id","title","content","outline" }
        return db.Select(strings.Join(fields, ", "))
    }).First(&shots)
    return
}

func GetBy(field string, value uint64) (shots Shots) {
    database.DB.Where(map[string]interface{}{field: value}).First(&shots)
    return
}

func GetByWhereMap(where map[string]interface{}) (shots Shots) {
	database.DB.Where(where).First(&shots)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (shots []Shots) {
	database.DB.Where(where).Find(&shots)
	return
}

func All() (shots []Shots) {
    database.DB.Find(&shots)
    return
}

func IsExist(field string, value uint64) bool {
    var count int64
    database.DB.Model(&Shots{}).Where(map[string]interface{}{field: value}).Count(&count)
    return count > 0
}

// Paginate 镜头表分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (shots []Shots, paging paginator.Paging) {
    // 构建带关联预加载的查询
    query := database.DB.Model(Shots{}).
        Preload("Projectss", func(db *gorm.DB) *gorm.DB {
            // 只加载关联表的必要字段
            fields := []string{"id","admin_id","serial_no","status","image" }
            return db.Select(strings.Join(fields, ", "))
    }).
        Preload("Scriptss", func(db *gorm.DB) *gorm.DB {
            // 只加载关联表的必要字段
            fields := []string{"id","project_id","title","content","outline" }
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
        &shots,
        app.V1URL(database.TableName(&Shots{})),
        perPage,
        )
    return
}