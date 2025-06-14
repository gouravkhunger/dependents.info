package cmd

import (
	"embed"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
	"dependents.info/internal/middleware"
	"dependents.info/internal/routes"
	"dependents.info/internal/service"
	"dependents.info/pkg/utils"
)

func Build(cfg *config.Config, static embed.FS, services *service.Services) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.Static(static))
	app.Use(middleware.Config(cfg))

	routes.Setup(app, services)

	return app
}
