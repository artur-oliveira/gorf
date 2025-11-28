package dto

import (
	"grf/core/dto"
	"time"
)

var _ dto.IPatchDTO = (*UserPatchDTO)(nil)

type UserCreateDTO struct {
	Username  string `json:"username" validate:"required,min=3,max=150"`
	Email     string `json:"email" validate:"required,email,max=254"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name,omitempty" validate:"max=150"`
	LastName  string `json:"last_name,omitempty" validate:"max=150"`
	IsStaff   *bool  `json:"is_staff,omitempty" validate:"omitempty,boolean"`
	IsActive  *bool  `json:"is_active,omitempty" validate:"omitempty,boolean"`
}

type UserUpdateDTO struct {
	Username    string `json:"username" validate:"required,min=3,max=150"`
	Email       string `json:"email" validate:"required,email,max=254"`
	FirstName   string `json:"first_name,omitempty" validate:"max=150"`
	LastName    string `json:"last_name,omitempty" validate:"max=150"`
	IsStaff     bool   `json:"is_staff" validate:"boolean"`
	IsActive    bool   `json:"is_active" validate:"boolean"`
	IsSuperuser bool   `json:"is_superuser" validate:"boolean"`
}

type UserPatchDTO struct {
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=150"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email,max=254"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,max=150"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,max=150"`
	IsStaff   *bool   `json:"is_staff,omitempty" validate:"omitempty,boolean"`
	IsActive  *bool   `json:"is_active,omitempty" validate:"omitempty,boolean"`
}

func (dto *UserPatchDTO) ToPatchMap() map[string]interface{} {
	updates := make(map[string]interface{})

	if dto.Username != nil {
		updates["username"] = *dto.Username
	}
	if dto.Email != nil {
		updates["email"] = *dto.Email
	}
	if dto.FirstName != nil {
		updates["first_name"] = *dto.FirstName
	}
	if dto.LastName != nil {
		updates["last_name"] = *dto.LastName
	}
	if dto.IsStaff != nil {
		updates["is_staff"] = *dto.IsStaff
	}
	if dto.IsActive != nil {
		updates["is_active"] = *dto.IsActive
	}

	return updates
}

func (dto *UserPatchDTO) IsEmpty() bool {
	return dto.Username == nil && dto.Email == nil && dto.FirstName == nil && dto.LastName == nil && dto.IsStaff == nil && dto.IsActive == nil
}

type UserResponseDTO struct {
	ID          uint64     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	IsActive    bool       `json:"is_active"`
	IsStaff     bool       `json:"is_staff"`
	IsSuperuser bool       `json:"is_superuser"`
	LastLogin   *time.Time `json:"last_login,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
