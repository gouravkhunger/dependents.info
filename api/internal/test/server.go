package test

import (
	"dependents-img/internal/config"
	"dependents-img/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func NewServer() *fiber.App {
	cfg := config.NewTestConfig()
	app := fiber.New()
	app.Use(middleware.Config(cfg))
	return app
}
