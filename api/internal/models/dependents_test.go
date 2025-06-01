package models

import (
	"testing"
)

func TestIngestRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     IngestRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: IngestRequest{
				Total: 8,
				Dependents: []Dependent{
					{
						Name:  "foo/bar",
						Image: "data:image/png;base64,r4nd0m=",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid request with different image type",
			req: IngestRequest{
				Total: 10,
				Dependents: []Dependent{
					{
						Name:  "owner/repo",
						Image: "data:image/jpeg;base64,r4nd0m=",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid name (wrong format)",
			req: IngestRequest{
				Total: -1,
				Dependents: []Dependent{
					{
						Name:  "foo-bar",
						Image: "data:image/png;base64,r4nd0m=",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid image (missing prefix)",
			req: IngestRequest{
				Total: 10,
				Dependents: []Dependent{
					{
						Name:  "foo/bar",
						Image: "r4nd0m=",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid totals count",
			req: IngestRequest{
				Total: -20,
				Dependents: []Dependent{
					{
						Name:  "foo/bar",
						Image: "r4nd0m=",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}
