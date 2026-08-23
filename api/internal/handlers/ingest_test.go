package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/env"
	"dependents.info/internal/models"
	"dependents.info/internal/test"
)

func TestIngestHandler_Ingest(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		expectedStatus int
		expectedMsg    string
		renderErr      error
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
				Total:      0,
				Dependents: []models.Dependent{{}},
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
						Image: "invalid_base64",
					},
				},
			},
			expectedMsg:    "Invalid JSON payload",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "render failure",
			requestBody: models.IngestRequest{
				Total: 10,
				Dependents: []models.Dependent{
					{Image: "data:image/png;base64,r4nd0m=="},
				},
			},
			renderErr:      errors.New("render failed"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedMsg:    "Failed to render SVG",
		},
		{
			name: "success",
			requestBody: models.IngestRequest{
				Total: 10,
				Dependents: []models.Dependent{
					{
						Image: "data:image/png;base64,r4nd0m==",
					},
				},
			},
			expectedStatus: fiber.StatusOK,
			expectedMsg:    "Dependents data ingested successfully",
		},
	}

	cfg := test.NewConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := test.NewMockStore()
			renderer := &test.MockRenderer{
				SVGResult: []byte("<svg>test</svg>"),
				SVGErr:    tt.renderErr,
			}

			app := test.NewServer(cfg)
			h := NewIngestHandler(nil, store, renderer)
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

func TestIngestHandler_Ingest_StoresData(t *testing.T) {
	cfg := test.NewConfig()
	store := test.NewMockStore()
	renderer := &test.MockRenderer{SVGResult: []byte("<svg>ok</svg>")}

	app := test.NewServer(cfg)
	h := NewIngestHandler(nil, store, renderer)
	app.Post("/:owner/:repo/ingest", h.Ingest)

	body := models.IngestRequest{
		Total:      42,
		Dependents: []models.Dependent{{Image: "data:image/png;base64,abc="}},
	}
	reqBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/owner/repo/ingest", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var total string
	if err := store.Get("total:owner/repo", &total); err != nil {
		t.Fatal("total not stored")
	}
	if total != "42" {
		t.Errorf("expected total '42', got %q", total)
	}

	var svg string
	if err := store.Get("svg:owner/repo", &svg); err != nil {
		t.Fatal("svg not stored")
	}
	if svg != "<svg>ok</svg>" {
		t.Errorf("expected svg '<svg>ok</svg>', got %q", svg)
	}
}

func TestIngestHandler_Ingest_ProductionAuth(t *testing.T) {
	cfg := test.NewConfig()
	cfg.Environment = env.EnvProduction

	t.Run("missing token", func(t *testing.T) {
		store := test.NewMockStore()
		renderer := &test.MockRenderer{SVGResult: []byte("<svg/>")}
		app := test.NewServer(cfg)
		h := NewIngestHandler(&test.MockOIDCVerifier{}, store, renderer)
		app.Post("/:owner/:repo/ingest", h.Ingest)

		body, _ := json.Marshal(models.IngestRequest{
			Total:      1,
			Dependents: []models.Dependent{{Image: "data:image/png;base64,x="}},
		})
		req := httptest.NewRequest("POST", "/owner/repo/ingest", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req, -1)
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		store := test.NewMockStore()
		renderer := &test.MockRenderer{SVGResult: []byte("<svg/>")}
		oidc := &test.MockOIDCVerifier{Err: errors.New("bad token")}
		app := test.NewServer(cfg)
		h := NewIngestHandler(oidc, store, renderer)
		app.Post("/:owner/:repo/ingest", h.Ingest)

		body, _ := json.Marshal(models.IngestRequest{
			Total:      1,
			Dependents: []models.Dependent{{Image: "data:image/png;base64,x="}},
		})
		req := httptest.NewRequest("POST", "/owner/repo/ingest", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer fake-token")
		resp, _ := app.Test(req, -1)
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Errorf("expected 401, got %d", resp.StatusCode)
		}
	})
}
