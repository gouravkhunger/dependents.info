package routes

import (
	"github.com/gofiber/fiber/v2"

	"dependents-img/internal/handlers"
	"dependents-img/internal/service"
)

func Setup(app *fiber.App, services *service.Services) {
	healthHandler := handlers.NewHealthHandler()
	ingestHandler := handlers.NewIngestHandler(services.GitHubOIDCService)

	app.Get("/health", healthHandler.Health)
	app.Post("/:owner/:repo/ingest", ingestHandler.Ingest)
}
