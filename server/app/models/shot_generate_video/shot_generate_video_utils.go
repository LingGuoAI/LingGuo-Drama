package shot_generate_video

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/pkg/app"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/paginator"
)

func Get(idStr string) (shotGenerateVideo ShotGenerateVideo) {
	database.DB.Where("id", idStr).First(&shotGenerateVideo)
	return
}

func GetBy(field string, value uint64) (shotGenerateVideo ShotGenerateVideo) {
	database.DB.Where(map[string]interface{}{field: value}).First(&shotGenerateVideo)
	return
}

func GetByWhereMap(where map[string]interface{}) (shotGenerateVideo ShotGenerateVideo) {
	database.DB.Where(where).First(&shotGenerateVideo)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (shotGenerateVideo []ShotGenerateVideo) {
	database.DB.Where(where).Find(&shotGenerateVideo)
	return
}

func All() (shotGenerateVideo []ShotGenerateVideo) {
	database.DB.Find(&shotGenerateVideo)
	return
}

func IsExist(field string, value uint64) bool {
	var count int64
	database.DB.Model(&ShotGenerateVideo{}).Where(map[string]interface{}{field: value}).Count(&count)
	return count > 0
}

// Paginate 场景项目分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (shotGenerateVideo []ShotGenerateVideo, paging paginator.Paging) {
	// 构建带关联预加载的查询
	query := database.DB.Model(ShotGenerateVideo{})

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
		&shotGenerateVideo,
		app.V1URL(database.TableName(&ShotGenerateVideo{})),
		perPage,
	)
	return
}
