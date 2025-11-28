package server

import (
	"grf/core/config"
	"grf/core/middleware"
	"grf/core/permission"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	FiberApp  *fiber.App
	Config    *config.Config
	DB        *gorm.DB
	Validator *validator.Validate

	I18nMw *middleware.I18NMiddleware

	Models []interface{}

	AllowAny                  permission.IPermission
	IsAuthenticatedOrReadOnly permission.IPermission
	IsAuthenticated           permission.IPermission
	IsAdmin                   permission.IPermission
}

func (a *App) Start() error {
	return a.FiberApp.Listen(":" + a.Config.ServerPort)
}
