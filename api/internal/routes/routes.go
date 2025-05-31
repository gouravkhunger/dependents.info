package routes

import (
	"dependents-img/internal/handlers"
	"dependents-img/internal/service"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, services *service.Services) {
	healthHandler := handlers.NewHealthHandler()
	ingestHandler := handlers.NewIngestHandler(services.GitHubOIDCService)

	app.Get("/health", healthHandler.Health)
	app.Post("/:owner/:repo/ingest", ingestHandler.Ingest)
}
