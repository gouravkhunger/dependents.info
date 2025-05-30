package app

import (
	"dependents-img/internal/api"
	"dependents-img/internal/config"
	"dependents-img/internal/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Start() error {
	cfg := config.New()
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(middleware.Logger())

	api.Setup(app)

	return app.Listen(fmt.Sprintf(":%s", cfg.Port))
}
