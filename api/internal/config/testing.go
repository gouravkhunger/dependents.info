package config

import (
	"dependents-img/internal/env"
)

func NewTestConfig() *Config {
	return &Config{
		Port:               "5000",
		Environment:        env.EnvDevelopment,
		GitHubOIDCAudience: "http://localhost:5000",
		GitHubOIDCIssuer:   "https://token.actions.githubusercontent.com",
	}
}
