package shot_frame_image

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/pkg/app"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/paginator"
)

func Get(idStr string) (shotFrameImages ShotFrameImages) {
	database.DB.Where("id", idStr).First(&shotFrameImages)
	return
}

func GetBy(field string, value uint64) (shotFrameImages ShotFrameImages) {
	database.DB.Where(map[string]interface{}{field: value}).First(&shotFrameImages)
	return
}

func GetByWhereMap(where map[string]interface{}) (shotFrameImages ShotFrameImages) {
	database.DB.Where(where).First(&shotFrameImages)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (shotFrameImages []ShotFrameImages) {
	database.DB.Where(where).Find(&shotFrameImages)
	return
}

func All() (shotFrameImages []ShotFrameImages) {
	database.DB.Find(&shotFrameImages)
	return
}

func IsExist(field string, value uint64) bool {
	var count int64
	database.DB.Model(&ShotFrameImages{}).Where(map[string]interface{}{field: value}).Count(&count)
	return count > 0
}

// Paginate 场景项目分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (shotFrameImages []ShotFrameImages, paging paginator.Paging) {
	// 构建带关联预加载的查询
	query := database.DB.Model(ShotFrameImages{})

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
		&shotFrameImages,
		app.V1URL(database.TableName(&ShotFrameImages{})),
		perPage,
	)
	return
}
