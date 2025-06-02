package middleware

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func Static(staticDir embed.FS) fiber.Handler {
	return filesystem.New(filesystem.Config{
		MaxAge:             86400,
		Root:               http.FS(staticDir),
		PathPrefix:         "/static",
		ContentTypeCharset: "text/html; charset=utf-8",
	})
}
