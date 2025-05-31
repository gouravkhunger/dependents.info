package test

import (
	"github.com/gofiber/fiber/v2"

	"dependents-img/internal/config"
	"dependents-img/internal/middleware"
)

func NewServer() *fiber.App {
	cfg := config.NewTestConfig()
	app := fiber.New()
	app.Use(middleware.Config(cfg))
	return app
}
