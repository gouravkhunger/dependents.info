package app

import (
	"dependents-img/internal/config"
	"dependents-img/internal/middleware"
	"dependents-img/internal/routes"
	"dependents-img/internal/service"
	"dependents-img/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Start() error {
	cfg := config.New()
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	app.Use(middleware.Config(cfg))
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	services := service.BuildAll(cfg)

	routes.Setup(app, services)

	return app.Listen(fmt.Sprintf(":%s", cfg.Port))
}
