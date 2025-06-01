package routes

import (
	"github.com/gofiber/fiber/v2"

	"dependents-img/internal/handlers"
	"dependents-img/internal/service"
)

func Setup(app *fiber.App, services *service.Services) {
	healthHandler := handlers.NewHealthHandler()
	imageHandler := handlers.NewImageHandler(services.DatabaseService)
	ingestHandler := handlers.NewIngestHandler(
		services.GitHubOIDCService,
		services.ImageService,
		services.DatabaseService,
	)

	app.Get("/health", healthHandler.Health)
	app.Post("/:owner/:repo/ingest", ingestHandler.Ingest)
	app.Get("/:owner/:repo/image.svg", imageHandler.SVGImage)
}
