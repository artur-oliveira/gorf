package service

import (
	"grf/core/dto"
	"grf/core/filterset"
	"grf/core/models"
	"grf/core/pagination"
	"grf/core/repository"
)

type IService[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet, ID comparable] interface {
	List(filter F, pagination pagination.IPagination[T]) (*pagination.Response[T], error)
	GetByID(id ID) (T, error)
	GetAllByID(ids []ID) ([]T, error)
	Create(dto C) (T, error)
	Update(id ID, dto U) (T, error)
	PartialUpdate(id ID, dto P) (T, error)
	Delete(id ID) error
}

type GenericService[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet, ID comparable] struct {
	Repo repository.IRepository[T, ID]

	MapCreateToModel func(dto C) T
	MapUpdateToModel func(dto U, model T) T
}

type Config[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet, ID comparable] struct {
	Repo repository.IRepository[T, ID]

	MapCreateToModel func(dto C) T
	MapUpdateToModel func(dto U, model T) T
}

func NewGenericService[T models.IModel, C any, U any, P dto.IPatchDTO, R any, F filterset.IFilterSet, ID comparable](
	config *Config[T, C, U, P, R, F, ID],
) IService[T, C, U, P, R, F, ID] {

	if config.Repo == nil || config.MapCreateToModel == nil || config.MapUpdateToModel == nil {
		panic("GenericService: Repo e Mappers de negócios (Create/Update) não podem ser nulos")
	}

	return &GenericService[T, C, U, P, R, F, ID]{
		Repo:             config.Repo,
		MapCreateToModel: config.MapCreateToModel,
		MapUpdateToModel: config.MapUpdateToModel,
	}
}

func (s *GenericService[T, C, U, P, R, F, ID]) List(
	filter F,
	pagination pagination.IPagination[T],
) (*pagination.Response[T], error) {
	return s.Repo.FindPaginated(filter, pagination)
}

func (s *GenericService[T, C, U, P, R, F, ID]) GetByID(id ID) (T, error) {
	return s.Repo.FindById(id)
}

func (s *GenericService[T, C, U, P, R, F, ID]) GetAllByID(ids []ID) ([]T, error) {
	return s.Repo.FindAllById(ids)
}

func (s *GenericService[T, C, U, P, R, F, ID]) Create(dto C) (T, error) {
	newRecord := s.MapCreateToModel(dto)

	err := s.Repo.Create(newRecord)

	return newRecord, err
}

func (s *GenericService[T, C, U, P, R, F, ID]) Update(id ID, dto U) (T, error) {
	record, err := s.Repo.FindById(id)
	if err != nil {
		return record, err
	}

	updatedRecord := s.MapUpdateToModel(dto, record)

	err = s.Repo.Update(updatedRecord)

	return updatedRecord, err
}

func (s *GenericService[T, C, U, P, R, F, ID]) PartialUpdate(id ID, dto P) (T, error) {
	record, err := s.Repo.FindById(id)
	if err != nil {
		return record, err
	}
	if dto.IsEmpty() {
		return record, nil
	}
	patchMap := dto.ToPatchMap()
	if len(patchMap) == 0 {
		return record, nil
	}

	err = s.Repo.PartialUpdate(record, patchMap)

	return record, err
}

func (s *GenericService[T, C, U, P, R, F, ID]) Delete(id ID) error {
	return s.Repo.Delete(id)
}
