package bootstrap

import (
	"grf/core/auth"
	"grf/core/config"
	"grf/core/database"
	"grf/core/exceptions"
	"grf/core/i18n"
	"grf/core/middleware"
	"grf/core/permission"
	"grf/core/routes"
	"grf/core/server"
	"grf/core/validator"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewApp(cfg config.Config, models []interface{}) (*server.App, error) {
	db, err := database.ConnectDB(&cfg)
	if err != nil {
		return nil, err
	}

	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: exceptions.GlobalErrorHandler,
		IdleTimeout:  time.Duration(cfg.ServerIdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ServerIdleTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerIdleTimeout) * time.Second,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(logger.New())

	jwtBackend := auth.NewJWTAuthBackend(db, &cfg)
	basicBackend := auth.NewBasicAuthBackend(db)

	i18nMw := middleware.NewI18NMiddleware(i18n.NewI18nService())

	isAuthenticated := permission.NewIsAuthenticated(jwtBackend, basicBackend)

	var bootstrapedApp = &server.App{
		FiberApp:  app,
		DB:        db,
		Validator: validator.GetValidator(),
		I18nMw:    i18nMw,
		Config:    &cfg,
		Models:    models,

		AllowAny:        &permission.AllowAny{},
		IsAuthenticated: isAuthenticated,
		IsAdmin:         &permission.IsAdmin{},
		IsAuthenticatedOrReadOnly: permission.NewOr(
			&permission.IsReadOnly{},
			isAuthenticated,
		),
	}

	i18nMw.UseMiddleWare(
		app,
	)
	database.RegisterMigrations(&database.MigrationOptions{
		DB:     db,
		Config: &cfg,
		Models: models,
	})
	permission.RegisterPermissions(&permission.Options{
		DB:     db,
		Models: models,
	})
	routes.RegisterRoutes(
		bootstrapedApp,
	)
	return bootstrapedApp, nil
}
