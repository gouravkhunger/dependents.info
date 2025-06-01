package test

import (
	"dependents.info/internal/config"
	"dependents.info/internal/env"
)

func NewConfig() *config.Config {
	return &config.Config{
		Port:               "5000",
		Environment:        env.EnvDevelopment,
		DatabasePath:       "/tmp/dependents-test",
		GitHubOIDCAudience: "http://localhost:5000",
		GitHubOIDCIssuer:   "https://token.actions.githubusercontent.com",
	}
}
