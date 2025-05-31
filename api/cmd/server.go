package app

import (
	"dependents-img/internal/api"
	"dependents-img/internal/config"
	"dependents-img/internal/middleware"
	"dependents-img/internal/service"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Start() error {
	cfg := config.New()
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	services := service.BuildAll(cfg)

	api.Setup(app, services)

	return app.Listen(fmt.Sprintf(":%s", cfg.Port))
}
