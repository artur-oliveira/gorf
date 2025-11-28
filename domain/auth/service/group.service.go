package service

import (
	"grf/core/exceptions"
	generic_repository "grf/core/repository"
	"grf/core/service"
	"grf/domain/auth/dto"
	"grf/domain/auth/filter"
	"grf/domain/auth/mapper"
	"grf/domain/auth/model"
	"grf/domain/auth/repository"

	"gorm.io/gorm"
)

type GroupService struct {
	service.IService[*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64]

	DB *gorm.DB
}

func NewGroupService(
	config *service.Config[*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64],
	db *gorm.DB,
) service.IService[*model.Group, *dto.GroupCreateDTO, *dto.GroupUpdateDTO, *dto.GroupPatchDTO, *dto.GroupResponseDTO, *filter.GroupFilterSet, uint64] {

	baseService := service.NewGenericService(config)
	return &GroupService{
		IService: baseService,
		DB:       db,
	}
}

func (s *GroupService) Create(dto *dto.GroupCreateDTO) (*model.Group, error) {
	newRecord := mapper.MapCreateToGroup(dto)

	err := s.doSaveOnTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(newRecord)
	}, newRecord, dto.PermissionIDs, generic_repository.SyncAlways)

	if err != nil {
		return nil, exceptions.NewInternal(err)
	}

	s.preload(newRecord)
	return newRecord, nil
}

func (s *GroupService) Update(id uint64, dto *dto.GroupUpdateDTO) (*model.Group, error) {
	record, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	updatedRecord := mapper.MapUpdateToGroup(dto, record)

	err = s.doSaveOnTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(updatedRecord)
	}, updatedRecord, dto.PermissionIDs, generic_repository.SyncAlways)

	if err != nil {
		return nil, exceptions.NewInternal(err)
	}

	s.preload(updatedRecord)
	return updatedRecord, nil
}

func (s *GroupService) PartialUpdate(id uint64, dto *dto.GroupPatchDTO) (*model.Group, error) {
	record, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	patchMap := dto.ToPatchMap()
	if len(patchMap) == 0 && dto.PermissionIDs == nil {
		return record, nil
	}

	err = s.doSaveOnTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(record).Updates(patchMap)
	}, record, dto.PermissionIDs, generic_repository.SyncIfProvided)

	if err != nil {
		return nil, err
	}

	s.preload(record)
	return record, nil
}

func (s *GroupService) doSaveOnTransaction(
	saveFn func(db *gorm.DB) *gorm.DB,
	record *model.Group,
	permissions []uint64,
	policy generic_repository.M2MSyncPolicy,
) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := saveFn(tx).Error; err != nil {
			return err
		}

		shouldSync := false
		if policy == generic_repository.SyncAlways {
			shouldSync = true
		} else if policy == generic_repository.SyncIfProvided && permissions != nil {
			shouldSync = true
		}

		if shouldSync {
			if err := s.syncPermissions(tx, record, permissions); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *GroupService) syncPermissions(
	tx *gorm.DB,
	group *model.Group,
	permissionIDs []uint64,
) error {
	permissionRepo := repository.NewPermissionRepository(tx)
	var permissions []*model.Permission
	var err error

	if len(permissionIDs) > 0 {
		permissions, err = permissionRepo.FindAllById(permissionIDs)
		if err != nil {
			return exceptions.NewBadRequest("error_query_permissions", err)
		}
	}

	if err := tx.Model(group).Association("Permissions").Replace(permissions); err != nil {
		return err
	}
	return nil
}

func (s *GroupService) preload(group *model.Group) {
	s.DB.Preload("Permissions").First(group, group.ID)
}
