package repository

import (
	"errors"
	"grf/core/filterset"
	"grf/core/models"
	"grf/core/pagination"

	"gorm.io/gorm"
)

type IRepository[T models.IModel, ID comparable] interface {
	FindPaginated(filter filterset.IFilterSet, pagination pagination.IPagination[T]) (*pagination.Response[T], error)
	FindById(id ID) (T, error)
	FindAllById(ids []ID) ([]T, error)
	Create(entity T) error
	CreateMany(entity []T) error
	Update(entity T) error
	PartialUpdate(entity T, updates map[string]interface{}) error
	Delete(id ID) error
}

type GenericRepository[T models.IModel, ID comparable] struct {
	DB       *gorm.DB
	NewModel func() T
}

type Config[T models.IModel, ID comparable] struct {
	DB       *gorm.DB
	NewModel func() T
}

func NewGenericRepository[T models.IModel, ID comparable](
	config *Config[T, ID],
) IRepository[T, ID] {
	if config.DB == nil || config.NewModel == nil {
		panic("GenericRepository: DB and NewModel cannot be nil")
	}
	return &GenericRepository[T, ID]{
		DB:       config.DB,
		NewModel: config.NewModel,
	}
}

func (r *GenericRepository[T, ID]) FindPaginated(
	filter filterset.IFilterSet,
	pagination pagination.IPagination[T],
) (*pagination.Response[T], error) {
	var results []T
	query := filter.Apply(r.DB.Model(&results))

	return pagination.Paginate(query)
}

func (r *GenericRepository[T, ID]) FindById(id ID) (T, error) {
	model := r.NewModel()
	if err := r.DB.First(model, id).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (r *GenericRepository[T, ID]) FindAllById(ids []ID) ([]T, error) {
	if ids == nil || len(ids) == 0 {
		return []T{}, nil
	}
	var results []T
	if err := r.DB.Where("id IN ?", ids).Find(&results).Error; err != nil {
		return nil, err
	}

	if len(results) != len(ids) {
		return nil, errors.New("wrong number of results")
	}
	return results, nil
}

func (r *GenericRepository[T, ID]) Create(entity T) error {
	return handleTx(r.DB.Create(entity))
}

func (r *GenericRepository[T, ID]) CreateMany(entities []T) error {
	return handleTx(r.DB.Create(entities))
}

func (r *GenericRepository[T, ID]) Update(entity T) error {
	return handleTx(r.DB.Save(entity))
}

func (r *GenericRepository[T, ID]) PartialUpdate(entity T, updates map[string]interface{}) error {
	return handleTx(r.DB.Model(entity).Updates(updates))
}

func (r *GenericRepository[T, ID]) Delete(id ID) error {
	record := r.NewModel()
	return handleTx(r.DB.Delete(record, id))
}

func handleTx(tx *gorm.DB) error {
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
