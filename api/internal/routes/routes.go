package routes

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/handlers"
	"dependents.info/internal/service"
)

func Setup(app *fiber.App, services *service.Services) {
	healthHandler := handlers.NewHealthHandler()
	imageHandler := handlers.NewImageHandler(services.DatabaseService)
	deleteHandler := handlers.NewDeleteHandler(services.DatabaseService)
	badgeHandler := handlers.NewBadgeHandler(
		services.DatabaseService,
		services.DependentsService,
	)
	sitemapHandler := handlers.NewSitemapHandler(
		services.DatabaseService,
		services.RenderService,
	)
	repoHandler := handlers.NewRepoHandler(
		services.DatabaseService,
		services.RenderService,
	)
	ingestHandler := handlers.NewIngestHandler(
		services.GitHubOIDCService,
		services.DatabaseService,
		services.RenderService,
	)

	app.Get("/health", healthHandler.Health)
	app.Get("/sitemap.xml", sitemapHandler.Sitemap)

	app.Get("/:owner/:repo", repoHandler.RepoPage)
	app.Delete("/:owner/:repo", deleteHandler.Delete)

	app.Get("/:owner/:repo/badge", badgeHandler.Badge)
	app.Get("/:owner/:repo/image", imageHandler.SVGImage)

	app.Get("/:owner/:repo/badge.svg", badgeHandler.Badge)
	app.Get("/:owner/:repo/image.svg", imageHandler.SVGImage)

	app.Post("/:owner/:repo/ingest", ingestHandler.Ingest)
}
