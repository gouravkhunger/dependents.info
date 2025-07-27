package handlers

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service/database"
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
			name:           "error",
			repo:           "owner/diff",
			password: 		 	"admin",
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "forbidden",
			repo:           "owner/repo",
			password: 		 	"wrongpassword",
			expectedStatus: fiber.StatusForbidden,
		},
		{
			name:           "success",
			repo:           "owner/repo",
			password: 		 	"admin",
			expectedStatus: fiber.StatusOK,
		},
	}

	cfg := test.NewConfig()
	dbService := database.NewBadgerService(cfg.DatabasePath)
	defer dbService.Close()

	dbService.Save("total:owner:repo", []byte("69420"))
	dbService.Save("svg:owner:repo", []byte("some svg string"))

	var data string
	dbService.Get("total:owner:repo", &data)
	log.Print(data)
	if data == "" {
		t.Fatalf("no totals count set for test repo")
	}
	dbService.Get("svg:owner:repo", &data)
	if data == "" {
		t.Fatalf("no svg image set for test repo")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := test.NewServer(cfg)
			req := httptest.NewRequest("DELETE", "/" + tt.repo, nil)
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
