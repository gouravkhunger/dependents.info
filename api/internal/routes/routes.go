package routes

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/handlers"
	"dependents.info/internal/service"
)

func Setup(app *fiber.App, services *service.Services) {
	healthHandler := handlers.NewHealthHandler()
	repoHandler := handlers.NewRepoHandler(services.RenderService)
	imageHandler := handlers.NewImageHandler(services.DatabaseService)
	badgeHandler := handlers.NewBadgeHandler(services.DatabaseService)
	ingestHandler := handlers.NewIngestHandler(
		services.GitHubOIDCService,
		services.DatabaseService,
		services.RenderService,
	)

	app.Get("/health", healthHandler.Health)
	app.Get("/:owner/:repo", repoHandler.RepoPage)
	app.Get("/:owner/:repo/badge", badgeHandler.Badge)
	app.Post("/:owner/:repo/ingest", ingestHandler.Ingest)
	app.Get("/:owner/:repo/image.svg", imageHandler.SVGImage)
}
