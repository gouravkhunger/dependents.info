package handlers

import (
	"dependents-img/internal/models"
	"dependents-img/internal/test"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestHealthHandler_Health(t *testing.T) {
	app := test.NewServer()
	h := NewHealthHandler()
	app.Get("/health", h.Health)

	tests := []struct {
		name           string
		expectedStatus int
		expectedMsg    string
	}{
		{
			expectedStatus: fiber.StatusOK,
			name:           "healthy response",
			expectedMsg:    "Service is healthy!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/health", nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			var apiResp models.APIResponse
			err = json.NewDecoder(resp.Body).Decode(&apiResp)
			if err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if !apiResp.Success {
				t.Errorf("expected Success true, got false")
			}
			if apiResp.Message != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, apiResp.Message)
			}

			// Check Data field
			dataBytes, _ := json.Marshal(apiResp.Data)
			var healthResp models.HealthResponse
			_ = json.Unmarshal(dataBytes, &healthResp)
			if healthResp.Status != "healthy" {
				t.Errorf("expected health status 'healthy', got %q", healthResp.Status)
			}
			if _, err = time.Parse(time.RFC3339, healthResp.Timestamp); err != nil {
				t.Errorf("timestamp is not valid RFC3339: %v", err)
			}
		})
	}
}

func TestHealthHandler_Health_MethodNotAllowed(t *testing.T) {
	app := test.NewServer()
	h := NewHealthHandler()
	app.Get("/health", h.Health)

	req := httptest.NewRequest("POST", "/health", strings.NewReader(""))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", fiber.StatusMethodNotAllowed, resp.StatusCode)
	}
}
