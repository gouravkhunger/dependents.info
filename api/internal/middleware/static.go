package middleware

import (
	"io/fs"
	"mime"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func Static(staticDir fs.FS) fiber.Handler {
	_ = mime.AddExtensionType(".webmanifest", "application/manifest+json")
	return filesystem.New(filesystem.Config{
		MaxAge:             86400,
		Root:               http.FS(staticDir),
		PathPrefix:         "/static",
		ContentTypeCharset: "utf-8",
	})
}
