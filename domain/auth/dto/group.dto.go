package dto

type GroupCreateDTO struct {
	Name          string   `json:"name" validate:"required,max=50"`
	PermissionIDs []uint64 `json:"permission_ids" validate:"omitempty,dive,gt=0"`
}

type GroupUpdateDTO struct {
	Name          string   `json:"name" validate:"required,max=50"`
	PermissionIDs []uint64 `json:"permission_ids" validate:"omitempty,dive,gt=0"`
}

type GroupResponseDTO struct {
	ID          uint64                  `json:"id"`
	Name        string                  `json:"name"`
	Permissions []PermissionResponseDTO `json:"permissions,omitempty"`
}

type GroupPatchDTO struct {
	Name          *string  `json:"name,omitempty" validate:"omitempty,max=50"`
	PermissionIDs []uint64 `json:"permission_ids,omitempty" validate:"omitempty,dive,gt=0"`
}

func (dto *GroupPatchDTO) ToPatchMap() map[string]interface{} {
	updates := make(map[string]interface{})
	if dto.Name != nil {
		updates["name"] = *dto.Name
	}
	return updates
}

func (dto *GroupPatchDTO) IsEmpty() bool {
	return dto.Name == nil
}
