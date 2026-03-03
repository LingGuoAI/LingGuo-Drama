package shot_video_merge

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/pkg/app"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/paginator"
)

func Get(idStr string) (shotVideoMerge ShotVideoMerge) {
	database.DB.Where("id", idStr).First(&shotVideoMerge)
	return
}

func GetBy(field string, value uint64) (shotVideoMerge ShotVideoMerge) {
	database.DB.Where(map[string]interface{}{field: value}).First(&shotVideoMerge)
	return
}

func GetByWhereMap(where map[string]interface{}) (shotVideoMerge ShotVideoMerge) {
	database.DB.Where(where).First(&shotVideoMerge)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (shotVideoMerge []ShotVideoMerge) {
	database.DB.Where(where).Find(&shotVideoMerge)
	return
}

func All() (shotVideoMerge []ShotVideoMerge) {
	database.DB.Find(&shotVideoMerge)
	return
}

func IsExist(field string, value uint64) bool {
	var count int64
	database.DB.Model(&ShotVideoMerge{}).Where(map[string]interface{}{field: value}).Count(&count)
	return count > 0
}

// Paginate 场景项目分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (shotVideoMerge []ShotVideoMerge, paging paginator.Paging) {
	// 构建带关联预加载的查询
	query := database.DB.Model(ShotVideoMerge{})

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
		&shotVideoMerge,
		app.V1URL(database.TableName(&ShotVideoMerge{})),
		perPage,
	)
	return
}
