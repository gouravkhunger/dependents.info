package utils

import (
	"testing"
)

func TestValidateRepository(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Valid cases
		{"a/b", true},
		{"owner/repo", true},
		{"owner/repo.git", true},
		{"abc-123/def_456", true},
		{"owner/repo.name", true},
		{"owner-1/repo-2.3_4", true},
		{"Owner-123/repo_name", true},
		// Invalid cases
		{"", false},
		{"/repo", false},
		{"owner", false},
		{"owner/", false},
		{"owner/repo/", false},
		{"owner@/repo", false},
		{"owner/repo!", false},
		{"owner//repo", false},
		{"owner/repo name", false},
		{"owner/repo$name", false},
		{"owner/repo/extra", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ValidateRepository(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateRepository(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "Valid Bearer token",
			header:    "Bearer abcdef123456",
			wantToken: "abcdef123456",
			wantErr:   false,
		},
		{
			name:      "Valid Bearer token with spaces",
			header:    "Bearer    tokenwithspaces",
			wantToken: "tokenwithspaces",
			wantErr:   false,
		},
		{
			name:      "No Bearer prefix",
			header:    "Token abcdef",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "Bearer only, no token",
			header:    "Bearer ",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "Empty header",
			header:    "",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "Bearer prefix in the middle",
			header:    "Token Bearer abcdef",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "Bearer prefix, short string",
			header:    "Bear",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "Bearer prefix, exact length",
			header:    "Bearer",
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := ExtractBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractBearerToken(%q) error = %v, wantErr %v", tt.header, err, tt.wantErr)
			}
			if token != tt.wantToken {
				t.Errorf("ExtractBearerToken(%q) = %q, want %q", tt.header, token, tt.wantToken)
			}
		})
	}
}
