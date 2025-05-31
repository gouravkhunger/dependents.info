package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents-img/internal/models"
	"dependents-img/internal/test"
)

func TestSendResponse(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		response models.APIResponse
	}{
		{
			name:   "Success with data",
			status: fiber.StatusOK,
			response: models.APIResponse{
				Success: true,
				Message: "OK",
				Data:    map[string]string{"foo": "bar"},
			},
		},
		{
			name:   "Success without data",
			status: fiber.StatusOK,
			response: models.APIResponse{
				Success: true,
				Message: "No Data",
			},
		},
		{
			name:   "Failure with error",
			status: fiber.StatusBadRequest,
			response: models.APIResponse{
				Success: false,
				Message: "Bad Request",
				Error:   "invalid input",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := test.NewServer()

			app.Get("/", func(c *fiber.Ctx) error {
				return SendResponse(c, tt.status, tt.response)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() error = %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			var got models.APIResponse
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("json.Unmarshal error = %v", err)
			}

			if got.Success != tt.response.Success || got.Message != tt.response.Message || got.Error != tt.response.Error {
				t.Errorf("expected %+v, got %+v", tt.response, got)
			}
		})
	}
}

func TestSendError(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		message  string
		err      error
		wantResp models.APIResponse
	}{
		{
			name:    "Error with error message",
			status:  fiber.StatusInternalServerError,
			message: "Internal Error",
			err:     errors.New("something went wrong"),
			wantResp: models.APIResponse{
				Success: false,
				Message: "Internal Error",
				Error:   "something went wrong",
			},
		},
		{
			name:    "Error without error message",
			status:  fiber.StatusBadRequest,
			message: "Bad Request",
			err:     nil,
			wantResp: models.APIResponse{
				Success: false,
				Message: "Bad Request",
				Error:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := test.NewServer()

			app.Get("/", func(c *fiber.Ctx) error {
				return SendError(c, tt.status, tt.message, tt.err)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() error = %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, resp.StatusCode)
			}

			body, _ := io.ReadAll(resp.Body)
			var got models.APIResponse
			if err := json.Unmarshal(body, &got); err != nil {
				t.Fatalf("json.Unmarshal error = %v", err)
			}

			if got.Success != tt.wantResp.Success || got.Message != tt.wantResp.Message || got.Error != tt.wantResp.Error {
				t.Errorf("expected %+v, got %+v", tt.wantResp, got)
			}
		})
	}
}
