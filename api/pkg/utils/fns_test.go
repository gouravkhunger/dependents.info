package utils

import (
	"strconv"
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

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{input: 500, expected: "500"},
		{input: 999, expected: "999"},
		{input: 1_000, expected: "1K"},
		{input: 1_200, expected: "1.2K"},
		{input: 12_660, expected: "12.7K"},
		{input: 15_420, expected: "15.4K"},
		{input: 19_890, expected: "19.9K"},
		{input: 29_990, expected: "30K"},
		{input: 100_623, expected: "101K"},
		{input: 818_020, expected: "818K"},
		{input: 948_563, expected: "949K"},
		{input: 999_490, expected: "999K"},
		{input: 999_563, expected: "1M"},
		{input: 999_999, expected: "1M"},
		{input: 1_000_000, expected: "1M"},
		{input: 1_250_000, expected: "1.3M"},
		{input: 9_999_999, expected: "10M"},
		{input: 12_500_000, expected: "12.5M"},
		{input: 999_999_599, expected: "1B"},
		{input: 999_999_999, expected: "1B"},
		{input: 1_000_000_000, expected: "1B"},
		{input: 2_345_000_000, expected: "2.3B"},
		{input: 2_350_000_000, expected: "2.4B"},
	}

	for _, tt := range tests {
		name := strconv.Itoa(tt.input)
		t.Run(name, func(t *testing.T) {
			result := FormatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("FormatNumber(%q) = %v; want %v", name, result, tt.expected)
			}
		})
	}
}

func TestToRoute(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"key", ""},
		{"r4nd0m:", ""},
		{":r4nd0m", ""},
		{"key:abcde", ""},
		{"key:owner/repo", "/owner/repo"},
		{"key:owner/repo:r4nd0mId=", "/owner/repo?id=r4nd0mId="},
		{"key:gouravkhunger/dependents.info", "/gouravkhunger/dependents.info"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToRoute(tt.input)
			if result != tt.expected {
				t.Errorf("ToURL(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSetParams(t *testing.T) {
	tests := []struct {
		input    string
		params   map[string]string
		expected string
	}{
		{
			input:    "https://example.com",
			params:   map[string]string{},
			expected: "https://example.com",
		},
		{
			input:    "https://example.com",
			params:   map[string]string{"foo": "bar"},
			expected: "https://example.com?foo=bar",
		},
		{
			input:    "https://example.com/path",
			params:   map[string]string{"a": "1", "b": "2"},
			expected: "https://example.com/path?a=1&b=2",
		},
		{
			input:    "https://example.com/path?x=5",
			params:   map[string]string{"y": "10", "z": ""},
			expected: "https://example.com/path?x=5&y=10",
		},
		{
			input:    "https://example.com/path?x=5",
			params:   map[string]string{},
			expected: "https://example.com/path?x=5",
		},
		{
			input:    "https://example.com/path",
			params:   map[string]string{"q": ""},
			expected: "https://example.com/path",
		},
		{
			input:    "https://example.com/path?foo=bar",
			params:   map[string]string{"foo": "baz"},
			expected: "https://example.com/path?foo=baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SetParams(tt.input, tt.params)
			if result != tt.expected {
				t.Errorf("SetParams(%q, %v) = %q; want %q", tt.input, tt.params, result, tt.expected)
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

func TestParseTotalDependents(t *testing.T) {
	html := `
		<div role="status" class="table-list-header-toggle states flex-auto pl-0">
			<a class="btn-link selected" href="/owner/repo/network/dependents?dependent_type=REPOSITORY&amp;package_id=someRandomID%3D">
				<svg aria-hidden="true" height="16" viewbox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-code-square">
					<path d="M0 1.75C0 .784.784 0 1.75 0h12.5C15.216 0 16 .784 16 1.75v12.5A1.75 1.75 0 0 1 14.25 16H1.75A1.75 1.75 0 0 1 0 14.25Zm1.75-.25a.25.25 0 0 0-.25.25v12.5c0 .138.112.25.25.25h12.5a.25.25 0 0 0 .25-.25V1.75a.25.25 0 0 0-.25-.25Zm7.47 3.97a.75.75 0 0 1 1.06 0l2 2a.75.75 0 0 1 0 1.06l-2 2a.749.749 0 0 1-1.275-.326.749.749 0 0 1 .215-.734L10.69 8 9.22 6.53a.75.75 0 0 1 0-1.06ZM6.78 6.53 5.31 8l1.47 1.47a.749.749 0 0 1-.326 1.275.749.749 0 0 1-.734-.215l-2-2a.75.75 0 0 1 0-1.06l2-2a.751.751 0 0 1 1.042.018.751.751 0 0 1 .018 1.042Z"></path>
				</svg>
				1,364
												Repositories
			</a>
			<a class="btn-link " href="/owner/repo/network/dependents?dependent_type=PACKAGE">
				<svg aria-hidden="true" height="16" viewbox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-package">
					<path d="m8.878.392 5.25 3.045c.54.314.872.89.872 1.514v6.098a1.75 1.75 0 0 1-.872 1.514l-5.25 3.045a1.75 1.75 0 0 1-1.756 0l-5.25-3.045A1.75 1.75 0 0 1 1 11.049V4.951c0-.624.332-1.201.872-1.514L7.122.392a1.75 1.75 0 0 1 1.756 0ZM7.875 1.69l-4.63 2.685L8 7.133l4.755-2.758-4.63-2.685a.248.248 0 0 0-.25 0ZM2.5 5.677v5.372c0 .09.047.171.125.216l4.625 2.683V8.432Zm6.25 8.271 4.625-2.683a.25.25 0 0 0 .125-.216V5.677L8.75 8.432Z"></path>
				</svg>
				709
												Packages
			</a>
			<details class="details-reset d-inline-block details-overlay js-dropdown-details position-relative">
				<summary aria-label="Warning" class="d-block px-1">
					<svg aria-hidden="true" height="16" viewbox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-info">
						<path d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8Zm8-6.5a6.5 6.5 0 1 0 0 13 6.5 6.5 0 0 0 0-13ZM6.5 7.75A.75.75 0 0 1 7.25 7h1a.75.75 0 0 1 .75.75v2.75h.25a.75.75 0 0 1 0 1.5h-2a.75.75 0 0 1 0-1.5h.25v-2h-.25a.75.75 0 0 1-.75-.75ZM8 6a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z"></path>
					</svg>
				</summary>
				<div class="Popover mt-2 right-0 mr-n2">
					<div class="Popover-message Popover-message--large Box color-shadow-large p-3 Popover-message--top-right ws-normal">
						These counts are approximate and may not exactly match the dependents shown below.
					</div>
				</div>
			</details>
		</div>
	`
	got, err := ParseTotalDependents(html, "owner/repo")
	if err != nil {
		t.Fatalf("ParseTotalDependents returned error: %v", err)
	}
	want := "1364"
	if got != want {
		t.Errorf("ParseTotalDependents() = %s, want %s", got, want)
	}
}
