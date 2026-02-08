package scenes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"spiritFruit/pkg/app"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/paginator"
	"strings"
)

func Get(idStr string) (scenes Scenes) {
	database.DB.Where("id", idStr).First(&scenes)
	return
}

func GetBy(field string, value uint64) (scenes Scenes) {
	database.DB.Where(map[string]interface{}{field: value}).First(&scenes)
	return
}

func GetByWhereMap(where map[string]interface{}) (scenes Scenes) {
	database.DB.Where(where).First(&scenes)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (scenes []Scenes) {
	database.DB.Where(where).Find(&scenes)
	return
}

func All() (scenes []Scenes) {
	database.DB.Find(&scenes)
	return
}

func IsExist(field string, value uint64) bool {
	var count int64
	database.DB.Model(&Scenes{}).Where(map[string]interface{}{field: value}).Count(&count)
	return count > 0
}

// Paginate 场景项目分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (scenes []Scenes, paging paginator.Paging) {
	// 构建带关联预加载的查询
	query := database.DB.Model(Scenes{}).
		Preload("Projectss", func(db *gorm.DB) *gorm.DB {
			// 只加载关联表的必要字段
			fields := []string{"id", "admin_id", "serial_no", "title", "description", "status", "image", "total_duration"}
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
		&scenes,
		app.V1URL(database.TableName(&Scenes{})),
		perPage,
	)
	return
}
