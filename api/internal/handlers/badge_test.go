package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service/database"
	"dependents.info/internal/service/github"
	"dependents.info/internal/service/render"
	"dependents.info/internal/test"
)

func TestBadgeHandler_Badge(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "valid request",
			url:            "/owner/repo/badge",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "non-existent badge",
			url:            "/invalid/repo/badge",
			expectedStatus: fiber.StatusNotFound,
		},
		{
			name:           "self",
			url:            "/self/repo/badge",
			expectedStatus: fiber.StatusOK,
		},
	}

	cfg := test.NewConfig()
	imageService := render.NewRenderService()
	dbService := database.NewBadgerService(cfg.DatabasePath)
	dependentsService := github.NewDependentsService(imageService)
	dbService.Save("total:owner/repo", []byte("69"))
	defer dbService.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := test.NewServer(cfg)
			h := NewBadgeHandler(dbService, dependentsService)
			app.Get("/self/repo/badge", h.SelfBadge)
			app.Get("/:owner/:repo/badge", h.Badge)

			req := httptest.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestBadgeHandler_Shields(t *testing.T) {
	cfg := test.NewConfig()
	cfg.DatabasePath = "/tmp/dependents-test-shields"
	imageService := render.NewRenderService()
	dbService := database.NewBadgerService(cfg.DatabasePath)
	dependentsService := github.NewDependentsService(imageService)
	dbService.Save("total:owner/repo", []byte("69"))
	dbService.Save("total:owner/repo:pkg1", []byte("1200"))
	defer dbService.Close()

	app := test.NewServer(cfg)
	h := NewBadgeHandler(dbService, dependentsService)
	app.Get("/:owner/:repo/shields.json", h.Shields)

	t.Run("cached total", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/owner/repo/shields.json", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("invalid json: %v", err)
		}
		if payload["schemaVersion"] != float64(1) {
			t.Errorf("schemaVersion = %v", payload["schemaVersion"])
		}
		if payload["label"] != "dependents" {
			t.Errorf("label = %v", payload["label"])
		}
		if payload["message"] != "69" {
			t.Errorf("message = %v", payload["message"])
		}
		if payload["color"] != "yellowgreen" {
			t.Errorf("color = %v", payload["color"])
		}
	})

	t.Run("package id and label override", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/owner/repo/shields.json?id=pkg1&label=used%20by", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("invalid json: %v", err)
		}
		if payload["label"] != "used by" {
			t.Errorf("label = %v", payload["label"])
		}
		if payload["message"] != "1.2K" {
			t.Errorf("message = %v", payload["message"])
		}
		if payload["color"] != "brightgreen" {
			t.Errorf("color = %v", payload["color"])
		}
	})

	t.Run("missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/missing/repo/shields.json", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
	})
}
