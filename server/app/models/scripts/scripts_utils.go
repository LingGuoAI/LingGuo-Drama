package scripts

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"spiritFruit/pkg/app"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/paginator"
	"strings"
)

func Get(idStr string) (scripts Scripts) {
	database.DB.Where("id", idStr).
		Preload("Projectss", func(db *gorm.DB) *gorm.DB {
			// 只加载关联表的必要字段
			fields := []string{"id", "admin_id", "serial_no", "title"}
			return db.Select(strings.Join(fields, ", "))
		}).Preload("Source").First(&scripts)
	return
}

func GetBy(field string, value uint64) (scripts Scripts) {
	database.DB.Where(map[string]interface{}{field: value}).First(&scripts)
	return
}

func GetByWhereMap(where map[string]interface{}) (scripts Scripts) {
	database.DB.Where(where).First(&scripts)
	return
}

func GetMapDataByWhereMap(where map[string]interface{}) (scripts []Scripts) {
	database.DB.Where(where).Find(&scripts)
	return
}

func All() (scripts []Scripts) {
	database.DB.Find(&scripts)
	return
}

func IsExist(field string, value uint64) bool {
	var count int64
	database.DB.Model(&Scripts{}).Where(map[string]interface{}{field: value}).Count(&count)
	return count > 0
}

// Paginate 剧本分页查询
func Paginate(c *gin.Context, perPage int, filters map[string]interface{}) (scripts []Scripts, paging paginator.Paging) {
	// 构建带关联预加载的查询
	query := database.DB.Model(Scripts{}).
		Select("scripts.*, (SELECT COUNT(*) FROM shots WHERE shots.script_id = scripts.id) as shots_count").
		Preload("Projectss", func(db *gorm.DB) *gorm.DB {
			// 只加载关联表的必要字段
			fields := []string{"id", "admin_id", "serial_no", "title"}
			return db.Select(strings.Join(fields, ", "))
		}).Preload("Source").
		Preload("ShotVideoMerges", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
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
		&scripts,
		app.V1URL(database.TableName(&Scripts{})),
		perPage,
	)
	return
}
