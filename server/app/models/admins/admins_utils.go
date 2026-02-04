package admins

import (
    "spiritFruit/pkg/app"
    "spiritFruit/pkg/database"
    "spiritFruit/pkg/paginator"
    "github.com/gin-gonic/gin"
)


func Get(idStr string) (admins Admins) {
    database.DB.Where("id", idStr).First(&admins)
    return
}

func GetBy(field string, value uint64) (admins Admins) {
    database.DB.Where(map[string]interface{}{field: value}).First(&admins)
    return
}

func GetByWhereMap(where map[string]interface{}) (admins Admins) {
	database.DB.Where(where).First(&admins)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (admins []Admins) {
	database.DB.Where(where).Find(&admins)
	return
}

func All() (admins []Admins) {
    database.DB.Find(&admins)
    return
}

func IsExist(field string, value uint64) bool {
    var count int64
    database.DB.Model(&Admins{}).Where(map[string]interface{}{field: value}).Count(&count)
    return count > 0
}

// Paginate 系统管理员分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (admins []Admins, paging paginator.Paging) {
    paging = paginator.Paginate(
        c,
        database.DB.Model(Admins{}),
        &admins,
        app.V1URL(database.TableName(&Admins{})),
        perPage,
        filters,
    )
    return
}