
package projects
import (
    "spiritFruit/app/models"
    "spiritFruit/pkg/database"
)

// Projects 结构体 短剧项目
type Projects struct {
    models.BaseModel
    AdminId  *uint64 `json:"adminId" form:"adminId" gorm:"default:1;column:admin_id;comment:归属用户ID(默认1);"`  //归属用户ID(默认1)
    SerialNo  *string `json:"serialNo" form:"serialNo" gorm:"column:serial_no;comment:业务流水号;size:64;"`  //业务流水号
    Title  *string `json:"title" form:"title" gorm:"column:title;comment:项目名称/短剧标题;size:255;"`  //项目名称/短剧标题
    Description  *string `json:"description" form:"description" gorm:"column:description;comment:项目简介;"`  //项目简介
    Status  *int8 `json:"status" form:"status" gorm:"default:0;column:status;comment:状态 0-Draft 1-Generating 2-Completed;"`  //状态
    Image  *string `json:"image" form:"image" gorm:"column:image;comment:封面图;size:1024;"`  //封面图
    TotalDuration  *uint64 `json:"totalDuration" form:"totalDuration" gorm:"default:0;column:total_duration;comment:总时长(秒);"`  //总时长(秒)
    Settings  *string `json:"settings" form:"settings" gorm:"column:settings;comment:生成配置快照;"`  //生成配置快照
    models.CommonTimestampsField
}


// TableName 短剧项目 Projects自定义表名 projects
func (Projects) TableName() string {
    return "projects"
}


// Create 创建短剧项目
func (projects *Projects) Create() {
    database.DB.Create(&projects)
}

// Save 保存短剧项目
func (projects *Projects) Save() (rowsAffected int64) {
    result := database.DB.Save(&projects)
    return result.RowsAffected
}

// Delete 删除短剧项目
func (projects *Projects) Delete() (rowsAffected int64) {
    result := database.DB.Delete(&projects)
    return result.RowsAffected
}

