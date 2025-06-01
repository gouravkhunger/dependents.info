package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/models"
	"dependents.info/internal/service/database"
	"dependents.info/internal/service/image"
	"dependents.info/internal/test"
)

func TestIngestHandler_Ingest(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "empty request body",
			requestBody:    models.IngestRequest{},
			expectedMsg:    "Invalid JSON payload",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid json",
			requestBody:    "{invalid json}",
			expectedMsg:    "Invalid JSON payload",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "invalid request validation",
			requestBody: models.IngestRequest{
				Total: 0,
				Dependents: []models.Dependent{
					{Name: "test"},
				},
			},
			expectedMsg:    "Invalid JSON payload",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "invalid image base64",
			requestBody: models.IngestRequest{
				Total: 0,
				Dependents: []models.Dependent{
					{
						Name:  "owner/repo",
						Image: "invalid_base64",
					},
				},
			},
			expectedMsg:    "Invalid JSON payload",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "success",
			requestBody: models.IngestRequest{
				Total: 10,
				Dependents: []models.Dependent{
					{
						Name:  "owner/repo",
						Image: "data:image/png;base64,r4nd0m==",
					},
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedMsg:    "Dependents data ingested successfully",
		},
	}

	cfg := test.NewConfig()
	imageService := image.NewImageService()
	dbService := database.NewBadgerService(cfg.DatabasePath)
	defer dbService.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := test.NewServer(cfg)
			h := NewIngestHandler(nil, imageService, dbService)
			app.Post("/ingest", h.Ingest)

			var reqBody []byte
			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/ingest", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
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
			if apiResp.Message != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, apiResp.Message)
			}
		})
	}
}
