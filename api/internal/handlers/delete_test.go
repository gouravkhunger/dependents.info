package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/test"
)

func TestDeleteHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		repo           string
		password       string
		expectedStatus int
	}{
		{
			name:           "wrong password",
			repo:           "owner/repo",
			password:       "wrongpassword",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "success",
			repo:           "owner/repo",
			password:       "admin",
			expectedStatus: fiber.StatusNoContent,
		},
	}

	cfg := test.NewConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := test.NewMockStore()
			store.Save("total:owner/repo", []byte("69420"))
			store.Save("svg:owner/repo", []byte("some svg string"))

			app := test.NewServer(cfg)
			h := NewDeleteHandler(store)
			app.Delete("/:owner/:repo", h.Delete)
			req := httptest.NewRequest("DELETE", "/"+tt.repo, nil)
			req.Header.Set("Authorization", "Bearer "+tt.password)
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

func TestDeleteHandler_Delete_RemovesKeys(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	store.Save("total:owner/repo", []byte("100"))
	store.Save("svg:owner/repo", []byte("<svg/>"))

	app := test.NewServer(cfg)
	h := NewDeleteHandler(store)
	app.Delete("/:owner/:repo", h.Delete)

	req := httptest.NewRequest("DELETE", "/owner/repo", nil)
	req.Header.Set("Authorization", "Bearer admin")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	var out string
	if err := store.Get("total:owner/repo", &out); err == nil {
		t.Error("expected total key to be deleted")
	}
	if err := store.Get("svg:owner/repo", &out); err == nil {
		t.Error("expected svg key to be deleted")
	}
}

func TestDeleteHandler_Delete_WithQueryId(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	store.Save("total:owner/repo:pkgId", []byte("50"))
	store.Save("svg:owner/repo:pkgId", []byte("<svg/>"))

	app := test.NewServer(cfg)
	h := NewDeleteHandler(store)
	app.Delete("/:owner/:repo", h.Delete)

	req := httptest.NewRequest("DELETE", "/owner/repo?id=pkgId", nil)
	req.Header.Set("Authorization", "Bearer admin")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	var out string
	if err := store.Get("total:owner/repo:pkgId", &out); err == nil {
		t.Error("expected total key to be deleted")
	}
}

func TestDeleteHandler_Delete_MissingAuth(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()

	app := test.NewServer(cfg)
	h := NewDeleteHandler(store)
	app.Delete("/:owner/:repo", h.Delete)

	req := httptest.NewRequest("DELETE", "/owner/repo", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}
