
package shots
import (
    "spiritFruit/app/models"
    "spiritFruit/pkg/database"
    "spiritFruit/app/models/projects"
    "spiritFruit/app/models/scripts"
)

// Shots 结构体 镜头表
type Shots struct {
    models.BaseModel
    ProjectId  *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;comment:所属项目ID, 外键约束(project_id) -> projects(id);"`  //所属项目ID
    ScriptId  *uint64 `json:"scriptId" form:"scriptId" gorm:"column:script_id;comment:所属剧本/分集ID, 外键约束(script_id) -> scripts(id);"`  //所属剧本/分集ID
    SequenceNo  *uint64 `json:"sequenceNo" form:"sequenceNo" gorm:"default:0;column:sequence_no;comment:镜头序号;"`  //镜头序号
    ShotType  *string `json:"shotType" form:"shotType" gorm:"column:shot_type;comment:景别: 全景/特写/中景;size:50;"`  //景别: 全景/特写/中景
    CameraMovement  *string `json:"cameraMovement" form:"cameraMovement" gorm:"column:camera_movement;comment:运镜: 推/拉/摇/移;size:50;"`  //运镜: 推/拉/摇/移
    Angle  *string `json:"angle" form:"angle" gorm:"column:angle;comment:视角: 俯视/平视;size:50;"`  //视角: 俯视/平视
    Dialogue  *string `json:"dialogue" form:"dialogue" gorm:"column:dialogue;comment:台词/旁白;"`  //台词/旁白
    VisualDesc  *string `json:"visualDesc" form:"visualDesc" gorm:"column:visual_desc;comment:画面描述;"`  //画面描述
    Atmosphere  *string `json:"atmosphere" form:"atmosphere" gorm:"column:atmosphere;comment:氛围/环境描述;"`  //氛围/环境描述
    ImagePrompt  *string `json:"imagePrompt" form:"imagePrompt" gorm:"column:image_prompt;comment:绘画Prompt;"`  //绘画Prompt
    VideoPrompt  *string `json:"videoPrompt" form:"videoPrompt" gorm:"column:video_prompt;comment:视频生成Prompt;"`  //视频生成Prompt
    AudioPrompt  *string `json:"audioPrompt" form:"audioPrompt" gorm:"column:audio_prompt;comment:音效/BGM提示词;"`  //音效/BGM提示词
    ImageUrl  *string `json:"imageUrl" form:"imageUrl" gorm:"column:image_url;comment:分镜图;size:1024;"`  //分镜图
    VideoUrl  *string `json:"videoUrl" form:"videoUrl" gorm:"column:video_url;comment:最终视频片段;size:1024;"`  //最终视频片段
    AudioUrl  *string `json:"audioUrl" form:"audioUrl" gorm:"column:audio_url;comment:配音/音效;size:1024;"`  //配音/音效
    DurationMs  *uint64 `json:"durationMs" form:"durationMs" gorm:"default:3000;column:duration_ms;comment:时长(毫秒, 原duration*1000);"`  //时长(毫秒, 原duration*1000)
    Status  *int8 `json:"status" form:"status" gorm:"default:0;column:status;comment:状态 0-Pending 1-Done 2-Fail;"`  //状态

    // 关联关系
    Projectss *projects.Projects `json:"projects,omitempty" gorm:"foreignKey:ProjectId;references:ID"` // 所属短剧项目
    Scriptss *scripts.Scripts `json:"scripts,omitempty" gorm:"foreignKey:ScriptId;references:ID"` // 所属剧本
    models.CommonTimestampsField
}


// TableName 镜头表 Shots自定义表名 shots
func (Shots) TableName() string {
    return "shots"
}


// Create 创建镜头表
func (shots *Shots) Create() {
    database.DB.Create(&shots)
}

// Save 保存镜头表
func (shots *Shots) Save() (rowsAffected int64) {
    result := database.DB.Save(&shots)
    return result.RowsAffected
}

// Delete 删除镜头表
func (shots *Shots) Delete() (rowsAffected int64) {
    result := database.DB.Delete(&shots)
    return result.RowsAffected
}

