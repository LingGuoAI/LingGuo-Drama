package source

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/pkg/app"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/paginator"
)

func Get(idStr string) (source Source) {
	database.DB.Where("id", idStr).First(&source)
	return
}

func GetBy(field string, value uint64) (source Source) {
	database.DB.Where(map[string]interface{}{field: value}).First(&source)
	return
}

func GetByWhereMap(where map[string]interface{}) (source Source) {
	database.DB.Where(where).First(&source)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (source []Source) {
	database.DB.Where(where).Find(&source)
	return
}

func All() (source []Source) {
	database.DB.Find(&source)
	return
}

func IsExist(field string, value uint64) bool {
	var count int64
	database.DB.Model(&Source{}).Where(map[string]interface{}{field: value}).Count(&count)
	return count > 0
}

// Paginate 场景项目分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (source []Source, paging paginator.Paging) {
	// 构建带关联预加载的查询
	query := database.DB.Model(Source{})

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
		&source,
		app.V1URL(database.TableName(&Source{})),
		perPage,
	)
	return
}
