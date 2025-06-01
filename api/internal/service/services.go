package service

import (
	"dependents.info/internal/config"
	"dependents.info/internal/service/database"
	"dependents.info/internal/service/github"
	"dependents.info/internal/service/image"
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
