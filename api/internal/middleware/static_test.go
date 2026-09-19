package middleware

import (
	"embed"
	"io/fs"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

//go:embed testdata
var testdata embed.FS

func TestStatic_ContentTypes(t *testing.T) {
	root, err := fs.Sub(testdata, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	app.Use(Static(root))

	tests := []struct {
		name        string
		url         string
		contentType string
	}{
		{
			name:        "llms.txt charset",
			url:         "/llms.txt",
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:        "webmanifest type",
			url:         "/site.webmanifest",
			contentType: "application/manifest+json; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != fiber.StatusOK {
				t.Fatalf("status %d", resp.StatusCode)
			}
			if ct := resp.Header.Get("Content-Type"); ct != tt.contentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.contentType)
			}
		})
	}
}
