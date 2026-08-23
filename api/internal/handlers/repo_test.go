package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/test"
)

func TestRepoHandler_RepoPage(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		storeData      map[string][]byte
		pageResult     []byte
		expectedStatus int
	}{
		{
			name:           "redirect when no data cached",
			url:            "/owner/repo",
			expectedStatus: fiber.StatusTemporaryRedirect,
		},
		{
			name:           "render page with total",
			url:            "/owner/repo",
			storeData:      map[string][]byte{"total:owner/repo": []byte("42")},
			pageResult:     []byte("<html>repo page</html>"),
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "render page with total and svg",
			url:  "/owner/repo",
			storeData: map[string][]byte{
				"total:owner/repo": []byte("42"),
				"svg:owner/repo":   []byte("<svg/>"),
			},
			pageResult:     []byte("<html>repo page</html>"),
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "with package id - redirect",
			url:            "/owner/repo?id=pkg1",
			expectedStatus: fiber.StatusTemporaryRedirect,
		},
		{
			name:           "with package id - render",
			url:            "/owner/repo?id=pkg1",
			storeData:      map[string][]byte{"total:owner/repo:pkg1": []byte("10")},
			pageResult:     []byte("<html>pkg page</html>"),
			expectedStatus: fiber.StatusOK,
		},
	}

	cfg := test.NewConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := test.NewMockStore()
			for k, v := range tt.storeData {
				store.Save(k, v)
			}

			renderer := &test.MockRenderer{PageResult: tt.pageResult}
			app := test.NewServer(cfg)
			h := NewRepoHandler(store, renderer)
			app.Get("/:owner/:repo", h.RepoPage)

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

func TestRepoHandler_Formats(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	store.Save("total:owner/repo", []byte("42"))
	store.Save("svg:owner/repo", []byte("<svg/>"))
	store.Save("total:owner/repo:pkg1", []byte("7"))

	app := test.NewServer(cfg)
	h := NewRepoHandler(store, &test.MockRenderer{PageResult: []byte("<html>ok</html>")})
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
		if !strings.Contains(text, "found **42 projects** (network dependents) using this github repository's default package!") {
			t.Errorf("missing wording: %s", text)
		}
		if !strings.Contains(text, "/owner/repo/badge") {
			t.Errorf("missing badge: %s", text)
		}
		if !strings.Contains(text, "/owner/repo/image") {
			t.Errorf("missing image: %s", text)
		}
		if !strings.Contains(text, "Made with [dependents.info](http://localhost:5000).") {
			t.Errorf("missing made with: %s", text)
		}
	})

	t.Run("markdown with package id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/owner/repo.md?id=pkg1", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
		body, _ := io.ReadAll(resp.Body)
		text := string(body)
		if !strings.Contains(text, "## package: pkg1") {
			t.Errorf("missing package heading: %s", text)
		}
		if !strings.Contains(text, "found **7 projects** (network dependents) using this github package!") {
			t.Errorf("missing package wording: %s", text)
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
