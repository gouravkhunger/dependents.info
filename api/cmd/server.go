package cmd

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
	"dependents.info/internal/middleware"
	"dependents.info/internal/routes"
	"dependents.info/internal/service"
	"dependents.info/pkg/utils"
)

func Build(cfg *config.Config, services *service.Services) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	app.Use(middleware.Config(cfg))
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.ETAG())

	routes.Setup(app, services)

	return app
}
