package service

import (
	"dependents-img/internal/config"
	"dependents-img/internal/service/github"
)

type Services struct {
	GitHubOIDCService *github.OIDCService
}

func BuildAll(cfg *config.Config) *Services {
	oidcService := github.NewOIDCService(cfg)

	return &Services{
		GitHubOIDCService: oidcService,
	}
}
