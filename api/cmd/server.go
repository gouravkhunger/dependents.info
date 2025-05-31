package cmd

import (
	"github.com/gofiber/fiber/v2"

	"dependents-img/internal/config"
	"dependents-img/internal/middleware"
	"dependents-img/internal/routes"
	"dependents-img/internal/service"
	"dependents-img/pkg/utils"
)

func Build(cfg *config.Config, services *service.Services) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	app.Use(middleware.Config(cfg))
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	routes.Setup(app, services)

	return app
}
