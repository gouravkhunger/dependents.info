package config

import (
	"context"
	"os"

	"dependents.info/internal/env"
)

type contextKey string

const ConfigContextKey = contextKey("config")

type Config struct {
	Port               string
	Password           string
	StylesFile         string
	DatabasePath       string
	GitHubOIDCAudience string
	GitHubOIDCIssuer   string
	Environment        env.Environment
}

func New() *Config {
	return &Config{
		Port:               getEnv("PORT", "5000"),
		Password:           getEnv("PASSWORD", "admin"),
		StylesFile:         getEnv("STYLES_FILE", "index.css"),
		DatabasePath:       getEnv("DATABASE_PATH", "/tmp/dependents"),
		Environment:        env.EnvFromString(getEnv("ENVIRONMENT", "development")),
		GitHubOIDCAudience: getEnv("GITHUB_OIDC_AUDIENCE", "https://dependents.info"),
		GitHubOIDCIssuer:   "https://token.actions.githubusercontent.com",
	}
}

func (c *Config) Host() string {
	if c.Environment == env.EnvProduction {
		return "https://dependents.info"
	} else {
		return "http://localhost:" + c.Port
	}
}

func FromContext(ctx context.Context) *Config {
	val := ctx.Value(ConfigContextKey)
	if cfg, ok := val.(*Config); ok {
		return cfg
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
