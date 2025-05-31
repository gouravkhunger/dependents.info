package config

import (
	"dependents-img/internal/env"
	"os"
)

type Config struct {
	Port               string
	GitHubOIDCAudience string
	GitHubOIDCIssuer   string
	Environment        env.Environment
}

func New() *Config {
	return &Config{
		Port:               getEnv("PORT", "5000"),
		Environment:        env.EnvFromString(getEnv("ENVIRONMENT", "development")),
		GitHubOIDCAudience: getEnv("GITHUB_OIDC_AUDIENCE", "https://dependents.info"),
		GitHubOIDCIssuer:   "https://token.actions.githubusercontent.com",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
