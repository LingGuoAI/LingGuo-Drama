package ai_config

import (
	"spiritFruit/pkg/app"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/paginator"

	"github.com/gin-gonic/gin"
)

func Get(idStr string) (config AiConfig) {
	database.DB.Where("id", idStr).First(&config)
	return
}

func GetBy(field string, value interface{}) (config AiConfig) {
	database.DB.Where(map[string]interface{}{field: value}).First(&config)
	return
}

func GetByWhereMap(where map[string]interface{}) (config AiConfig) {
	database.DB.Where(where).First(&config)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (configs []AiConfig) {
	database.DB.Where(where).Order("priority DESC, id DESC").Find(&configs)
	return
}

func All() (configs []AiConfig) {
	database.DB.Find(&configs)
	return
}

func IsExist(field string, value interface{}) bool {
	var count int64
	database.DB.Model(&AiConfig{}).Where(map[string]interface{}{field: value}).Count(&count)
	return count > 0
}

// Paginate 分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (configs []AiConfig, paging paginator.Paging) {
	paging = paginator.Paginate(
		c,
		// 按优先级倒序，再按ID倒序
		database.DB.Model(AiConfig{}).Order("priority DESC, id DESC"),
		&configs,
		app.V1URL(database.TableName(&AiConfig{})),
		perPage,
		filters,
	)
	return
}
