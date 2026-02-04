package admins
import (
	"spiritFruit/pkg/hash"
	"gorm.io/gorm"
)

// func (admins *Admins) BeforeSave(tx *gorm.DB) (err error) {}
// func (admins *Admins) BeforeCreate(tx *gorm.DB) (err error) {}
// func (admins *Admins) AfterCreate(tx *gorm.DB) (err error) {}
// func (admins *Admins) BeforeUpdate(tx *gorm.DB) (err error) {}
// func (admins *Admins) AfterUpdate(tx *gorm.DB) (err error) {}
// func (admins *Admins) AfterSave(tx *gorm.DB) (err error) {}
// func (admins *Admins) BeforeDelete(tx *gorm.DB) (err error) {}
// func (admins *Admins) AfterDelete(tx *gorm.DB) (err error) {}
// func (admins *Admins) AfterFind(tx *gorm.DB) (err error) {}
// BeforeSave GORM 的模型钩子，在创建和更新模型前调用
func (adminModel *Admins) BeforeSave(tx *gorm.DB) (err error) {
	if adminModel.Password != nil && *adminModel.Password != "" {
		if !hash.BcryptIsHashed(*adminModel.Password) {
			hashedPassword := hash.BcryptHash(*adminModel.Password)
			adminModel.Password = &hashedPassword
		}
	}
	return
}

