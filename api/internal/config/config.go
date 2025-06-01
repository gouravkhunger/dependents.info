package config

import (
	"context"
	"os"

	"dependents-img/internal/env"
)

type contextKey string

const ConfigContextKey = contextKey("config")

type Config struct {
	Port               string
	DatabasePath       string
	GitHubOIDCAudience string
	GitHubOIDCIssuer   string
	Environment        env.Environment
}

func New() *Config {
	return &Config{
		Port:               getEnv("PORT", "5000"),
		DatabasePath:       getEnv("DATABASE_PATH", "/tmp/dependents"),
		Environment:        env.EnvFromString(getEnv("ENVIRONMENT", "development")),
		GitHubOIDCAudience: getEnv("GITHUB_OIDC_AUDIENCE", "https://dependents.info"),
		GitHubOIDCIssuer:   "https://token.actions.githubusercontent.com",
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
