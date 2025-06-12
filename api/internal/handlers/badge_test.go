package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service/database"
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
	}

	cfg := test.NewConfig()
	dbService := database.NewBadgerService(cfg.DatabasePath)
	dbService.Save("total:owner/repo", []byte("42"))
	defer dbService.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := test.NewServer(cfg)
			h := NewBadgeHandler(dbService)
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
