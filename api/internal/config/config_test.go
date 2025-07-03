package config

import (
	"os"
	"testing"

	"dependents.info/internal/env"
)

func TestNew_WithEnvVars(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("STYLES_FILE", "test.css")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("DATABASE_PATH", "/tmp/testdb")
	os.Setenv("GITHUB_OIDC_AUDIENCE", "https://custom-audience.com")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("STYLES_FILE")
	defer os.Unsetenv("ENVIRONMENT")
	defer os.Unsetenv("DATABASE_PATH")
	defer os.Unsetenv("GITHUB_OIDC_AUDIENCE")

	cfg := New()

	if cfg.Port != "8080" {
		t.Errorf("expected Port to be '8080', got '%s'", cfg.Port)
	}
	if cfg.StylesFile != "test.css" {
		t.Errorf("expected StylesFile to be 'test.css', got '%s'", cfg.StylesFile)
	}
	if cfg.DatabasePath != "/tmp/testdb" {
		t.Errorf("expected DatabasePath to be '/tmp/testdb', got '%s'", cfg.DatabasePath)
	}
	if cfg.Environment != env.EnvProduction {
		t.Errorf("expected Environment to be 'production', got '%s'", cfg.Environment)
	}
	if cfg.GitHubOIDCAudience != "https://custom-audience.com" {
		t.Errorf("expected GitHubOIDCAudience to be 'https://custom-audience.com', got '%s'", cfg.GitHubOIDCAudience)
	}
	if cfg.Host() != "https://dependents.info" {
		t.Errorf("expected Host() to be '%s', got '%s'", "https://dependents.info", cfg.Host())
	}
}

func TestNew_WithDefaults(t *testing.T) {
	os.Unsetenv("PORT")
	os.Unsetenv("STYLES_FILE")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("DATABASE_PATH")
	os.Unsetenv("GITHUB_OIDC_AUDIENCE")

	cfg := New()

	if cfg.Port != "5000" {
		t.Errorf("expected default Port to be '5000', got '%s'", cfg.Port)
	}
	if cfg.StylesFile != "index.css" {
		t.Errorf("expected default StylesFile to be 'index.css', got '%s'", cfg.StylesFile)
	}
	if cfg.DatabasePath != "/tmp/dependents" {
		t.Errorf("expected default DatabasePath to be '/tmp/dependents', got '%s'", cfg.DatabasePath)
	}
	if cfg.Environment != env.EnvDevelopment {
		t.Errorf("expected default Environment to be 'development', got '%s'", cfg.Environment)
	}
	if cfg.GitHubOIDCAudience != "https://dependents.info" {
		t.Errorf("expected default GitHubOIDCAudience to be 'https://dependents.info', got '%s'", cfg.GitHubOIDCAudience)
	}
	if cfg.Host() != "http://localhost:5000" {
		t.Errorf("expected Host() to be '%s', got '%s'", "http://localhost:5000", cfg.Host())
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
