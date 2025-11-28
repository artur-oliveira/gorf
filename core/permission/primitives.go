package permission

import (
	"errors"
	"fmt"
	"grf/core/auth"
	"grf/core/exceptions"
	"grf/core/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AllowAny struct{}

func (p *AllowAny) Check(_ *fiber.Ctx) error {
	return nil
}

type IsReadOnly struct{}

func (p *IsReadOnly) Check(c *fiber.Ctx) error {
	if IsReadOnlyMethod(c.Method()) {
		return nil
	}
	return exceptions.NewForbidden("permission_denied_read_only", nil)
}

type IsAuthenticated struct {
	Backends []auth.IAuthBackend
}

func NewIsAuthenticated(backends ...auth.IAuthBackend) *IsAuthenticated {
	if len(backends) == 0 {
		panic("IsAuthenticated requer pelo menos um auth.IAuthBackend")
	}
	return &IsAuthenticated{Backends: backends}
}

func (p *IsAuthenticated) Check(c *fiber.Ctx) error {
	if c.Locals("user") != nil {
		return nil
	}

	for _, backend := range p.Backends {
		user, err := backend.Authenticate(c)

		if err != nil {
			if errors.Is(err, auth.ErrCannotAuthenticate) {
				continue
			}
			return err
		}

		if user != nil {
			if !user.Active() {
				return exceptions.NewUnauthorized("inactive_user", nil)
			}
			c.Locals("user", user)
			return nil
		}
	}
	return exceptions.NewUnauthorized("auth_invalid_or_not_provided", nil)
}

type IsAdmin struct {
}

func (p *IsAdmin) Check(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		return err
	}
	if !user.Admin() {
		return exceptions.NewForbidden("admin_required", nil) // 403 Forbidden
	}
	return nil
}

type ModelPermissions struct {
	DB    *gorm.DB
	Model models.IModel
}

func NewModelPermissions(db *gorm.DB, model models.IModel) *ModelPermissions {
	return &ModelPermissions{
		DB:    db,
		Model: model,
	}
}

func (p *ModelPermissions) Check(c *fiber.Ctx) error {
	user, err := GetUser(c)
	if err != nil {
		return err
	}

	action, err := getActionForContext(c)
	if err != nil {
		return exceptions.NewInternal(err)
	}
	if action == "" {
		return nil
	}
	if !user.HasPerm(p.DB, p.Model.ModuleName(), action) {
		permKey := p.Model.ModuleName() + "." + action
		return exceptions.NewForbidden(fmt.Sprintf("error_auth_permission_denied %s", permKey), nil)
	}
	return nil
}
