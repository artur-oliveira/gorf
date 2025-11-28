package repository

import (
	"grf/core/repository"
	"grf/domain/auth/model"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	repository.IRepository[*model.Permission, uint64]

	DB *gorm.DB
}

func NewPermissionRepository(
	db *gorm.DB,
) *PermissionRepository {
	return &PermissionRepository{
		IRepository: repository.NewGenericRepository(&repository.Config[*model.Permission, uint64]{
			DB: db,
			NewModel: func() *model.Permission {
				return new(model.Permission)
			},
		}),
		DB: db,
	}
}
