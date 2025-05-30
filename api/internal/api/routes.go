package api

import (
	"dependents-img/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	healthHandler := handlers.NewHealthHandler()
	app.Get("/health", healthHandler.Health)
}
