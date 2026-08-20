package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service/database"
	"dependents.info/internal/service/render"
	"dependents.info/internal/test"
)

func TestRepoHandler_Formats(t *testing.T) {
	cfg := test.NewConfig()
	cfg.DatabasePath = "/tmp/dependents-test-repo-formats"
	dbService := database.NewBadgerService(cfg.DatabasePath)
	dbService.Save("total:owner/repo", []byte("42"))
	dbService.Save("svg:owner/repo", []byte("<svg/>"))
	dbService.Save("total:owner/repo:pkg1", []byte("7"))
	defer dbService.Close()

	app := test.NewServer(cfg)
	h := NewRepoHandler(dbService, render.NewRenderService())
	app.Get("/:owner/:repo", h.RepoPage)

	t.Run("json", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/owner/repo.json", nil)
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
		if payload["owner"] != "owner" || payload["repo"] != "repo" {
			t.Errorf("owner/repo = %v/%v", payload["owner"], payload["repo"])
		}
		if payload["total"] != float64(42) {
			t.Errorf("total = %v", payload["total"])
		}
		if payload["has_image"] != true {
			t.Errorf("has_image = %v", payload["has_image"])
		}
		urls, _ := payload["urls"].(map[string]any)
		if urls["json"] != "http://localhost:5000/owner/repo.json" {
			t.Errorf("urls.json = %v", urls["json"])
		}
		if urls["shields"] != "http://localhost:5000/owner/repo/shields.json" {
			t.Errorf("urls.shields = %v", urls["shields"])
		}
	})

	t.Run("markdown", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/owner/repo.md", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		text := string(body)
		if !strings.Contains(text, "# owner/repo") {
			t.Errorf("missing heading: %s", text)
		}
		if !strings.Contains(text, "42") {
			t.Errorf("missing total: %s", text)
		}
		if !strings.Contains(text, "/owner/repo/badge") {
			t.Errorf("missing badge: %s", text)
		}
		if !strings.Contains(text, "/owner/repo/image") {
			t.Errorf("missing image: %s", text)
		}
	})

	t.Run("json with package id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/owner/repo.json?id=pkg1", nil)
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
		if payload["id"] != "pkg1" {
			t.Errorf("id = %v", payload["id"])
		}
		if payload["total"] != float64(7) {
			t.Errorf("total = %v", payload["total"])
		}
		if payload["has_image"] != false {
			t.Errorf("has_image = %v", payload["has_image"])
		}
	})

	t.Run("json missing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/missing/repo.json", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusNotFound {
			t.Errorf("expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("html still works", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/owner/repo", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})
}
