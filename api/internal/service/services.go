package service

import (
	"dependents-img/internal/config"
	"dependents-img/internal/service/database"
	"dependents-img/internal/service/github"
	"dependents-img/internal/service/image"
)

type Services struct {
	GitHubOIDCService *github.OIDCService
	DatabaseService   *database.BadgerService
	ImageService      *image.ImageService
}

func BuildAll(cfg *config.Config) *Services {
	imageService := image.NewImageService()
	oidcService := github.NewOIDCService(cfg)
	dbService := database.NewBadgerService(cfg.DatabasePath)

	return &Services{
		GitHubOIDCService: oidcService,
		DatabaseService:   dbService,
		ImageService:      imageService,
	}
}
