package service

import (
	"dependents.info/internal/config"
	"dependents.info/internal/service/database"
	"dependents.info/internal/service/github"
	"dependents.info/internal/service/render"
)

type Services struct {
	GitHubOIDCService *github.OIDCService
	DatabaseService   *database.BadgerService
	RenderService     *render.RenderService
}

func BuildAll(cfg *config.Config) *Services {
	imageService := render.NewRenderService()
	oidcService := github.NewOIDCService(cfg)
	dbService := database.NewBadgerService(cfg.DatabasePath)

	return &Services{
		GitHubOIDCService: oidcService,
		DatabaseService:   dbService,
		RenderService:     imageService,
	}
}
