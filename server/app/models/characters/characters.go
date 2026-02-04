
package characters
import (
    "spiritFruit/app/models"
    "spiritFruit/pkg/database"
    "spiritFruit/app/models/projects"
)

// Characters 结构体 角色
type Characters struct {
    models.BaseModel
    ProjectId  *uint64 `json:"projectId" form:"projectId" gorm:"column:project_id;comment:所属项目ID, 外键约束(project_id) -> projects(id);"`  //所属项目ID
    Name  *string `json:"name" form:"name" gorm:"column:name;comment:角色名;size:100;"`  //角色名
    RoleType  *string `json:"roleType" form:"roleType" gorm:"default:main;column:role_type;comment:角色类型: main/supporting/minor;size:50;"`  //角色类型: main/supporting/minor
    Gender  *string `json:"gender" form:"gender" gorm:"column:gender;comment:性别(需从appearance解析或留空);size:20;"`  //性别(需从appearance解析或留空)
    AgeGroup  *string `json:"ageGroup" form:"ageGroup" gorm:"column:age_group;comment:年龄段;size:50;"`  //年龄段
    Personality  *string `json:"personality" form:"personality" gorm:"column:personality;comment:性格描述;"`  //性格描述
    AppearanceDesc  *string `json:"appearanceDesc" form:"appearanceDesc" gorm:"column:appearance_desc;comment:外貌长文本描述(原appearance);"`  //外貌长文本描述(原appearance)
    VisualPrompt  *string `json:"visualPrompt" form:"visualPrompt" gorm:"column:visual_prompt;comment:AI绘画专用Prompt(从appearance提取);"`  //AI绘画专用Prompt(从appearance提取)
    AvatarUrl  *string `json:"avatarUrl" form:"avatarUrl" gorm:"column:avatar_url;comment:头像/定妆照;size:1024;"`  //头像/定妆照
    VoiceId  *string `json:"voiceId" form:"voiceId" gorm:"column:voice_id;comment:TTS音色ID;size:100;"`  //TTS音色ID

    // 关联关系
    Projectss *projects.Projects `json:"projects,omitempty" gorm:"foreignKey:ProjectId;references:ID"` // 所属短剧项目
    models.CommonTimestampsField
}


// TableName 角色 Characters自定义表名 characters
func (Characters) TableName() string {
    return "characters"
}


// Create 创建角色
func (characters *Characters) Create() {
    database.DB.Create(&characters)
}

// Save 保存角色
func (characters *Characters) Save() (rowsAffected int64) {
    result := database.DB.Save(&characters)
    return result.RowsAffected
}

// Delete 删除角色
func (characters *Characters) Delete() (rowsAffected int64) {
    result := database.DB.Delete(&characters)
    return result.RowsAffected
}

