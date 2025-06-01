package test

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
	"dependents.info/internal/middleware"
)

func NewServer(cfg *config.Config) *fiber.App {
	app := fiber.New()
	app.Use(middleware.Config(cfg))
	return app
}
