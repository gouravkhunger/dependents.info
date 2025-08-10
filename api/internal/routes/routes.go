package routes

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/handlers"
)

func Setup(app *fiber.App, handlers *handlers.Handlers) {
	app.Get("/health", handlers.HealthHandler.Health)
	app.Get("/sitemap.xml", handlers.SitemapHandler.Sitemap)

	app.Get("/:owner/:repo", handlers.RepoHandler.RepoPage)
	app.Delete("/:owner/:repo", handlers.DeleteHandler.Delete)

	app.Get("/:owner/:repo/badge", handlers.BadgeHandler.Badge)
	app.Get("/:owner/:repo/image", handlers.ImageHandler.SVGImage)

	app.Get("/:owner/:repo/badge.svg", handlers.BadgeHandler.Badge)
	app.Get("/:owner/:repo/image.svg", handlers.ImageHandler.SVGImage)

	app.Post("/:owner/:repo/ingest", handlers.IngestHandler.Ingest)
}
