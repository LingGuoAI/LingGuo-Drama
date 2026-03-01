package shot_video_merge

import (
	"spiritFruit/app/models"
	"spiritFruit/pkg/database"
)

// ShotVideoMerge 分镜视频合并表模型
type ShotVideoMerge struct {
	models.BaseModel
	ProjectId *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;comment:项目ID;"`
	ScriptId  *uint64 `json:"scriptId" form:"scriptId" gorm:"column:script_id;comment:剧本ID;"`
	Title     *string `json:"title" form:"title" gorm:"column:title;size:255;comment:标题;"`
	MergedUrl *string `json:"mergedUrl" form:"mergedUrl" gorm:"column:merged_url;size:255;comment:合成后的视频地址;"`
	Duration  *int    `json:"duration" form:"duration" gorm:"column:duration;comment:视频时长;"`
	TaskId    *uint64 `json:"taskId" form:"taskId" gorm:"column:task_id;comment:绑定的异步任务ID;"`
	ErrorMsg  *string `json:"errorMsg" form:"errorMsg" gorm:"column:error_msg;type:text;comment:错误信息;"`

	models.CommonTimestampsField
}

// TableName 表名
func (ShotVideoMerge) TableName() string {
	return "shot_video_merge"
}

// Create 创建记录
func (shotVideoMerge *ShotVideoMerge) Create() {
	database.DB.Create(&shotVideoMerge)
}

// Save 保存记录
func (shotVideoMerge *ShotVideoMerge) Save() (rowsAffected int64) {
	result := database.DB.Save(&shotVideoMerge)
	return result.RowsAffected
}

// Delete 删除记录
func (shotVideoMerge *ShotVideoMerge) Delete() (rowsAffected int64) {
	result := database.DB.Delete(&shotVideoMerge)
	return result.RowsAffected
}
