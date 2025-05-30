package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAPIResponseSerialization(t *testing.T) {
	tests := []struct {
		name      string
		input     APIResponse
		wantJSON  string
		expectErr bool
	}{
		{
			name: "success serialization",
			input: APIResponse{
				Error:   "",
				Success: true,
				Message: "ok",
				Data:    map[string]any{"foo": "bar"},
			},
			wantJSON:  `{"success":true,"message":"ok","data":{"foo":"bar"}}`,
			expectErr: false,
		},
		{
			name: "wrong serialization",
			input: APIResponse{
				Error:   "",
				Success: false,
				Message: "fail",
				Data:    map[string]any{},
			},
			wantJSON:  `{"success":true,"message":"ok","data":{}}`,
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			//print err
			fmt.Printf("Error: %v\n", err)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			if !tt.expectErr && string(got) != tt.wantJSON {
				t.Errorf("json.Marshal() = %s, want %s", got, tt.wantJSON)
			}
		})
	}
}

func TestHealthResponseSerialization(t *testing.T) {
	tests := []struct {
		name      string
		input     HealthResponse
		wantJSON  string
		expectErr bool
	}{
		{
			name: "success serialization",
			input: HealthResponse{
				Status:    "healthy",
				Timestamp: "2024-06-01T12:00:00Z",
			},
			wantJSON:  `{"status":"healthy","timestamp":"2024-06-01T12:00:00Z"}`,
			expectErr: false,
		},
		{
			name: "wrong serialization",
			input: HealthResponse{
				Status:    "bad",
				Timestamp: "now",
			},
			wantJSON:  `{"status":"good","timestamp":"now"}`,
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}
			if !tt.expectErr && string(got) != tt.wantJSON {
				t.Errorf("json.Marshal() = %s, want %s", got, tt.wantJSON)
			}
		})
	}
}
