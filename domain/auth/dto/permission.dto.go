package dto

import (
	"grf/core/dto"
)

var _ dto.IPatchDTO = (*PermissionPatchDTO)(nil)

type PermissionCreateDTO struct {
	Module      string `json:"module" validate:"required,max=100"`
	Action      string `json:"action" validate:"required,max=100"`
	Description string `json:"description" validate:"required,max=255"`
}

type PermissionUpdateDTO struct {
	Description string `json:"description" validate:"required,max=255"`
}

type PermissionPatchDTO struct {
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
}

func (dto *PermissionPatchDTO) IsEmpty() bool {
	return dto.Description == nil
}

func (dto *PermissionPatchDTO) ToPatchMap() map[string]interface{} {
	updates := make(map[string]interface{})
	if dto.Description != nil {
		updates["description"] = *dto.Description
	}
	return updates
}

type PermissionResponseDTO struct {
	ID          uint64 `json:"id"`
	Module      string `json:"module"`
	Action      string `json:"action"`
	Description string `json:"description"`
}
