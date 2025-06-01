package test

import (
	"github.com/gofiber/fiber/v2"

	"dependents-img/internal/config"
	"dependents-img/internal/middleware"
)

func NewServer(cfg *config.Config) *fiber.App {
	app := fiber.New()
	app.Use(middleware.Config(cfg))
	return app
}
