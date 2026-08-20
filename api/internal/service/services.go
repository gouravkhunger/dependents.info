package service

import (
	"context"
	"time"

	"dependents.info/internal/config"
	"dependents.info/internal/models"
	"dependents.info/internal/service/database"
	"dependents.info/internal/service/github"
	"dependents.info/internal/service/render"
)

type Store interface {
	Get(key string, out *string) error
	Save(key string, data []byte) error
	SaveWithTTL(key string, data []byte, ttl time.Duration) error
	Delete(key string) error
	IterateKeys(callback func(key string))
	Close() error
	Sync() error
}

type Renderer interface {
	RenderSVG(d models.IngestRequest) ([]byte, error)
	RenderPage(d models.RepoPage) ([]byte, error)
	RenderSitemap(data []string) ([]byte, error)
}

type OIDCVerifier interface {
	VerifyToken(ctx context.Context, rawToken string, expectedRepo string) error
}

type DependentsTasker interface {
	NewTask(repo string, id string, kind string, callback func(total int, svg []byte)) error
}

type Services struct {
	OIDCVerifier      OIDCVerifier
	DependentsService DependentsTasker
	Store             Store
	Renderer          Renderer
}

func BuildAll(cfg *config.Config) *Services {
	imageService := render.NewRenderService()
	oidcService := github.NewOIDCService(cfg)
	dbService := database.NewBadgerService(cfg.DatabasePath)
	dependentsService := github.NewDependentsService(imageService)

	return &Services{
		OIDCVerifier:      oidcService,
		Store:             dbService,
		Renderer:          imageService,
		DependentsService: dependentsService,
	}
}
