package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Password    string `gorm:"size:128;not null"`
	LastLogin   *time.Time
	IsSuperuser bool   `gorm:"default:false"`
	Username    string `gorm:"size:150;uniqueIndex;not null"`
	FirstName   string `gorm:"size:150"`
	LastName    string `gorm:"size:150"`
	Email       string `gorm:"size:254;uniqueIndex;not null"`
	IsStaff     bool   `gorm:"default:false"`
	IsActive    bool   `gorm:"default:true"`

	Groups          []*Group      `gorm:"many2many:auth_user_groups;"`
	UserPermissions []*Permission `gorm:"many2many:auth_user_permissions;"`
}

func (u *User) TableName() string { return "auth_user" }

func (u *User) ModuleName() string { return "user" }

func (u *User) Active() bool { return u.IsActive }

func (u *User) Admin() bool { return u.IsSuperuser || u.IsStaff }

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) HasPerm(
	db *gorm.DB,
	module string,
	action string,
) bool {
	if u.IsActive && u.IsSuperuser {
		return true
	}
	if !u.IsActive {
		return false
	}

	var permID uint64
	err := db.Model(&Permission{}).
		Select("id").
		Where("module = ? AND action = ?", module, action).
		First(&permID).Error

	if err != nil {
		return false
	}

	var directPermissionCount int64
	err = db.Table("auth_user_permissions").
		Where("user_id = ? AND permission_id = ?", u.ID, permID).
		Count(&directPermissionCount).Error

	if err == nil && directPermissionCount > 0 {
		return true // Encontrada permissão direta
	}

	var groupPermissionCount int64
	err = db.Table("auth_user_groups").
		Select("auth_user_groups.group_id").
		Joins("INNER JOIN auth_group_permissions ON auth_user_groups.group_id = auth_group_permissions.group_id").
		Where("auth_user_groups.user_id = ?", u.ID).
		Where("auth_group_permissions.permission_id = ?", permID).
		Count(&groupPermissionCount).Error

	return err == nil && groupPermissionCount > 0
}
