package permission

import (
	basemodels "grf/core/models"
	"grf/core/repository"
	"grf/domain/auth/model"

	"gorm.io/gorm"
)

type Options struct {
	DB     *gorm.DB
	Models []interface{}
}

func RegisterPermissions(options *Options) {
	repo := repository.NewGenericRepository[*model.Permission, int64](&repository.Config[*model.Permission, int64]{
		DB:       options.DB,
		NewModel: func() *model.Permission { return new(model.Permission) },
	})
	var perms []*model.Permission
	for _, dstOpt := range options.Models {
		perms = append(perms, generateModelPermissions(dstOpt)...)
	}
	err := repo.CreateMany(perms)
	if err != nil {
		panic(err)
	}
}

func generateModelPermissions(dst interface{}) []*model.Permission {
	var permissions []*model.Permission
	newDst, ok := dst.(basemodels.IModel)
	if ok {
		permissions = append(permissions, &model.Permission{
			Module:      newDst.ModuleName(),
			Action:      basemodels.ListAction,
			Description: "Permission to list " + newDst.TableName() + " table records.",
		})
		permissions = append(permissions, &model.Permission{
			Module:      newDst.ModuleName(),
			Action:      basemodels.DetailAction,
			Description: "Permission to detail " + newDst.TableName() + " table record.",
		})
		permissions = append(permissions, &model.Permission{
			Module:      newDst.ModuleName(),
			Action:      basemodels.CreateAction,
			Description: "Permission to create " + newDst.TableName() + " record.",
		})
		permissions = append(permissions, &model.Permission{
			Module:      newDst.ModuleName(),
			Action:      basemodels.UpdateAction,
			Description: "Permission to update " + newDst.TableName() + " record.",
		})
		permissions = append(permissions, &model.Permission{
			Module:      newDst.ModuleName(),
			Action:      basemodels.PartialUpdateAction,
			Description: "Permission to partial update " + newDst.TableName() + " record.",
		})
		permissions = append(permissions, &model.Permission{
			Module:      newDst.ModuleName(),
			Action:      basemodels.DeleteAction,
			Description: "Permission to delete " + newDst.TableName() + " record.",
		})
	}
	return permissions

}
