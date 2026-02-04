
package admins
import (
    "spiritFruit/app/models"
    "spiritFruit/pkg/database"
    "spiritFruit/pkg/hash"
)

// Admins 结构体 系统管理员
type Admins struct {
    models.BaseModel
    Username  *string `json:"username" form:"username" gorm:"column:username;comment:用户名;size:120;"`  //用户名
    Mobile  *string `json:"mobile" form:"mobile" gorm:"column:mobile;comment:手机号;size:11;"`  //手机号
    Password  *string `json:"password" form:"password" gorm:"column:password;comment:密码;size:64;"`  //密码
    Email  *string `json:"email" form:"email" gorm:"column:email;comment:邮箱;size:80;"`  //邮箱
    AuthorityId  *uint64 `json:"authorityId" form:"authorityId" gorm:"default:888;column:authority_id;comment:用户角色ID;size:20;"`  //用户角色ID
    models.CommonTimestampsField
}


// TableName 系统管理员 Admins自定义表名 admins
func (Admins) TableName() string {
    return "admins"
}


// Create 创建系统管理员
func (admins *Admins) Create() {
    database.DB.Create(&admins)
}

// Save 保存系统管理员
func (admins *Admins) Save() (rowsAffected int64) {
    result := database.DB.Save(&admins)
    return result.RowsAffected
}

// Delete 删除系统管理员
func (admins *Admins) Delete() (rowsAffected int64) {
    result := database.DB.Delete(&admins)
    return result.RowsAffected
}
// GetByMulti 通过 手机号/用户名 来获取系统管理员
func GetByMulti(loginID string) (adminsModel Admins) {
    database.DB.
        Where("mobile = ?", loginID).
        Or("username = ?", loginID).
        First(&adminsModel)
    return
}

// ComparePassword 密码是否正确
func (adminsModel *Admins) ComparePassword(_password string) bool {
    return hash.BcryptCheck(_password, *adminsModel.Password)
}

