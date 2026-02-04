package projects

import (
    "spiritFruit/pkg/app"
    "spiritFruit/pkg/database"
    "spiritFruit/pkg/paginator"
    "github.com/gin-gonic/gin"
)


func Get(idStr string) (projects Projects) {
    database.DB.Where("id", idStr).First(&projects)
    return
}

func GetBy(field string, value uint64) (projects Projects) {
    database.DB.Where(map[string]interface{}{field: value}).First(&projects)
    return
}

func GetByWhereMap(where map[string]interface{}) (projects Projects) {
	database.DB.Where(where).First(&projects)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (projects []Projects) {
	database.DB.Where(where).Find(&projects)
	return
}

func All() (projects []Projects) {
    database.DB.Find(&projects)
    return
}

func IsExist(field string, value uint64) bool {
    var count int64
    database.DB.Model(&Projects{}).Where(map[string]interface{}{field: value}).Count(&count)
    return count > 0
}

// Paginate 短剧项目分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (projects []Projects, paging paginator.Paging) {
    paging = paginator.Paginate(
        c,
        database.DB.Model(Projects{}),
        &projects,
        app.V1URL(database.TableName(&Projects{})),
        perPage,
        filters,
    )
    return
}