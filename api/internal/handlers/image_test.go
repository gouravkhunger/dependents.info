package handlers

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/test"
)

func TestImageHandler_SVGImage(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		storeData      map[string][]byte
		taskFn         func(string, string, string, func(int, []byte)) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "cached svg returned",
			url:            "/owner/repo/image",
			storeData:      map[string][]byte{"svg:owner/repo": []byte("<svg>cached</svg>")},
			expectedStatus: fiber.StatusOK,
			expectedBody:   "<svg>cached</svg>",
		},
		{
			name: "fallback via NewTask succeeds",
			url:  "/owner/repo/image",
			taskFn: func(repo, id, kind string, cb func(int, []byte)) error {
				if cb != nil {
					cb(10, []byte("<svg>fresh</svg>"))
				}
				return nil
			},
			expectedStatus: fiber.StatusOK,
			expectedBody:   "<svg>fresh</svg>",
		},
		{
			name: "fallback via NewTask fails",
			url:  "/owner/repo/image",
			taskFn: func(_, _, _ string, _ func(int, []byte)) error {
				return errors.New("fetch failed")
			},
			expectedStatus: fiber.StatusNotFound,
		},
		{
			name:           "with package id",
			url:            "/owner/repo/image?id=pkg1",
			storeData:      map[string][]byte{"svg:owner/repo:pkg1": []byte("<svg>pkg</svg>")},
			expectedStatus: fiber.StatusOK,
			expectedBody:   "<svg>pkg</svg>",
		},
	}

	cfg := test.NewConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := test.NewMockStore()
			for k, v := range tt.storeData {
				store.Save(k, v)
			}

			depService := &test.MockDependentsTasker{NewTaskFn: tt.taskFn}
			app := test.NewServer(cfg)
			h := NewImageHandler(store, depService)
			app.Get("/:owner/:repo/image", h.SVGImage)

			req := httptest.NewRequest("GET", tt.url, nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
			if tt.expectedBody != "" {
				body, _ := io.ReadAll(resp.Body)
				if string(body) != tt.expectedBody {
					t.Errorf("expected body %q, got %q", tt.expectedBody, string(body))
				}
			}
		})
	}
}