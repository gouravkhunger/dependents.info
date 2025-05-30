package config

import (
	"dependents-img/internal/env"
	"os"
	"testing"
)

func TestNew_WithEnvVars(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("GITHUB_OIDC_AUDIENCE", "https://custom-audience.com")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("ENVIRONMENT")
	defer os.Unsetenv("GITHUB_OIDC_AUDIENCE")

	cfg := New()

	if cfg.Port != "8080" {
		t.Errorf("expected Port to be '8080', got '%s'", cfg.Port)
	}
	if cfg.Environment != env.EnvProduction {
		t.Errorf("expected Environment to be 'production', got '%s'", cfg.Environment)
	}
	if cfg.GitHubOIDCAudience != "https://custom-audience.com" {
		t.Errorf("expected GitHubOIDCAudience to be 'https://custom-audience.com', got '%s'", cfg.GitHubOIDCAudience)
	}
}

func TestNew_WithDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("GITHUB_OIDC_AUDIENCE")

	cfg := New()

	if cfg.Port != "5000" {
		t.Errorf("expected default Port to be '5000', got '%s'", cfg.Port)
	}
	if cfg.Environment != env.EnvDevelopment {
		t.Errorf("expected default Environment to be 'development', got '%s'", cfg.Environment)
	}
	if cfg.GitHubOIDCAudience != "https://dependents.info" {
		t.Errorf("expected default GitHubOIDCAudience to be 'https://dependents.info', got '%s'", cfg.GitHubOIDCAudience)
	}
}

func Test_getEnv(t *testing.T) {
	key := "TEST_KEY"
	defaultVal := "default"
	os.Unsetenv(key)
	if val := getEnv(key, defaultVal); val != defaultVal {
		t.Errorf("expected default value '%s', got '%s'", defaultVal, val)
	}

	os.Setenv(key, "set-value")
	defer os.Unsetenv(key)
	if val := getEnv(key, defaultVal); val != "set-value" {
		t.Errorf("expected env value 'set-value', got '%s'", val)
	}
}
