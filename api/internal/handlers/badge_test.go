package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/test"
)

func TestBadgeHandler_Badge(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		storeData      map[string][]byte
		taskErr        error
	}{
		{
			name:           "cached total returns badge",
			url:            "/owner/repo/badge",
			expectedStatus: fiber.StatusOK,
			storeData:      map[string][]byte{"total:owner/repo": []byte("69")},
		},
		{
			name:           "fallback via NewTask succeeds",
			url:            "/owner/repo/badge",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "fallback via NewTask fails",
			url:            "/owner/repo/badge",
			expectedStatus: fiber.StatusNotFound,
			taskErr:        errors.New("fetch failed"),
		},
	}

	cfg := test.NewConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := test.NewMockStore()
			for k, v := range tt.storeData {
				store.Save(k, v)
			}

			depService := &test.MockDependentsTasker{
				NewTaskFn: func(_ context.Context, repo, id, kind string, callback func(int, []byte)) error {
					if tt.taskErr != nil {
						return tt.taskErr
					}
					if callback != nil {
						callback(100, nil)
					}
					return nil
				},
			}

			app := test.NewServer(cfg)
			h := NewBadgeHandler(store, depService)
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
	store := test.NewMockStore()
	store.Save("total:owner/repo", []byte("69"))
	store.Save("total:owner/repo:pkg1", []byte("1200"))

	app := test.NewServer(cfg)
	h := NewBadgeHandler(store, &test.MockDependentsTasker{})
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

func TestBadgeHandler_SelfBadge(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	store.Save("total:owner1/repo1", []byte("10"))
	store.Save("svg:owner1/repo1", []byte("<svg/>"))
	store.Save("total:owner2/repo2", []byte("20"))

	app := test.NewServer(cfg)
	h := NewBadgeHandler(store, &test.MockDependentsTasker{})
	app.Get("/self/badge", h.SelfBadge)

	req := httptest.NewRequest("GET", "/self/badge", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestBadgeHandler_Badge_WithQueryId(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	store.Save("total:owner/repo:pkg123", []byte("42"))

	app := test.NewServer(cfg)
	h := NewBadgeHandler(store, &test.MockDependentsTasker{})
	app.Get("/:owner/:repo/badge", h.Badge)

	req := httptest.NewRequest("GET", "/owner/repo/badge?id=pkg123", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestColorFunction(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0", "red"},
		{"-1", "red"},
		{"5", "yellow"},
		{"50", "yellowgreen"},
		{"500", "green"},
		{"5000", "brightgreen"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := color(tt.input)
			if got != tt.want {
				t.Errorf("color(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
