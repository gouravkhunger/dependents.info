package service

import (
	"dependents-img/internal/config"
	"dependents-img/internal/service/database"
	"dependents-img/internal/service/github"
)

type Services struct {
	GitHubOIDCService *github.OIDCService
	DatabaseService   *database.BadgerService
}

func BuildAll(cfg *config.Config) *Services {
	oidcService := github.NewOIDCService(cfg)
	dbService := database.NewBadgerService(cfg.DatabasePath)

	return &Services{
		GitHubOIDCService: oidcService,
		DatabaseService:   dbService,
	}
}
