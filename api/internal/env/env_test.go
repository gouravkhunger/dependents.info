package env

import (
	"testing"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Environment
	}{
		{"empty input", "", EnvDevelopment},
		{"random", "r4nd0m", EnvDevelopment},
		{"production", "production", EnvProduction},
		{"development", "development", EnvDevelopment},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnvFromString(tt.input); got != tt.want {
				t.Errorf("EnvFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
