package handlers

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/test"
)

func TestSitemapHandler_Sitemap(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	store.Save("total:owner/repo", []byte("10"))
	store.Save("svg:owner/repo", []byte("<svg/>"))

	renderer := &test.MockRenderer{
		SitemapResult: []byte(`<?xml version="1.0"?><urlset><url><loc>http://localhost:5000/owner/repo</loc></url></urlset>`),
	}

	app := test.NewServer(cfg)
	h := NewSitemapHandler(store, renderer)
	app.Get("/sitemap.xml", h.Sitemap)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "owner/repo") {
		t.Error("expected sitemap to contain owner/repo")
	}
}

func TestSitemapHandler_Sitemap_Empty(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	renderer := &test.MockRenderer{SitemapResult: []byte(`<?xml version="1.0"?><urlset></urlset>`)}

	app := test.NewServer(cfg)
	h := NewSitemapHandler(store, renderer)
	app.Get("/sitemap.xml", h.Sitemap)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestSitemapHandler_Sitemap_RenderError(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	renderer := &test.MockRenderer{SitemapErr: errors.New("template error")}

	app := test.NewServer(cfg)
	h := NewSitemapHandler(store, renderer)
	app.Get("/sitemap.xml", h.Sitemap)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}