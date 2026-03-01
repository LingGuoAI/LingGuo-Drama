package shots

import (
	"spiritFruit/pkg/app"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/paginator"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Get(idStr string) (shots Shots) {
	database.DB.Where("id", idStr).
		Preload("Projects", func(db *gorm.DB) *gorm.DB {
			// 只加载关联表的必要字段
			fields := []string{"id", "admin_id", "serial_no", "status", "image"}
			return db.Select(strings.Join(fields, ", "))
		}).
		Preload("Scenes", func(db *gorm.DB) *gorm.DB {
			fields := []string{"id", "project_id", "name", "location", "visual_prompt"}
			return db.Select(strings.Join(fields, ", "))
		}).
		Preload("Scripts", func(db *gorm.DB) *gorm.DB {
			// 只加载关联表的必要字段
			fields := []string{"id", "project_id", "title", "content", "outline"}
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
		Preload("Projects", func(db *gorm.DB) *gorm.DB {
			// 只加载关联表的必要字段
			fields := []string{"id", "admin_id", "serial_no", "status", "image"}
			return db.Select(strings.Join(fields, ", "))
		}).
		Preload("Scenes", func(db *gorm.DB) *gorm.DB {
			fields := []string{"id", "project_id", "name", "location", "visual_prompt"}
			return db.Select(strings.Join(fields, ", "))
		}).
		Preload("Scripts", func(db *gorm.DB) *gorm.DB {
			// 只加载关联表的必要字段
			fields := []string{"id", "project_id", "title", "content", "outline"}
			return db.Select(strings.Join(fields, ", "))
		}).
		Preload("Characters", func(db *gorm.DB) *gorm.DB {
			fields := []string{"id", "name", "avatar_url", "visual_prompt"}
			return db.Select(strings.Join(fields, ", "))
		}).
		Preload("FrameImages", func(db *gorm.DB) *gorm.DB {
			fields := []string{"id", "shot_id", "frame_type", "image_type", "image_url", "created_at"}
			// 推荐按照创建时间倒序，这样前端拿到的第一张就是最新生成的
			return db.Select(strings.Join(fields, ", ")).Order("created_at DESC")
		}).
		Preload("Props", func(db *gorm.DB) *gorm.DB {
			fields := []string{"id", "name", "image_url"}
			return db.Select(strings.Join(fields, ", "))
		}).Preload("FramePrompts", func(db *gorm.DB) *gorm.DB {
		fields := []string{"id", "shot_id", "frame_type", "prompt", "description"}
		return db.Select(strings.Join(fields, ", "))
	}).Preload("GenerateVideos", func(db *gorm.DB) *gorm.DB {
		fields := []string{"id", "shot_id", "video_url", "created_at"}
		return db.Select(strings.Join(fields, ", ")).Order("created_at DESC")
	}).First(&shots)

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
