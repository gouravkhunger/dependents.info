package handlers

import "dependents.info/internal/service"

type Handlers struct {
	RepoHandler    *RepoHandler
	ImageHandler   *ImageHandler
	BadgeHandler   *BadgeHandler
	IngestHandler  *IngestHandler
	HealthHandler  *HealthHandler
	DeleteHandler  *DeleteHandler
	SitemapHandler *SitemapHandler
}

func BuildAll(services *service.Services) *Handlers {
	healthHandler := NewHealthHandler()
	deleteHandler := NewDeleteHandler(services.DatabaseService)
	imageHandler := NewImageHandler(services.DatabaseService, services.DependentsService)
	badgeHandler := NewBadgeHandler(
		services.DatabaseService,
		services.DependentsService,
	)
	sitemapHandler := NewSitemapHandler(
		services.DatabaseService,
		services.RenderService,
	)
	repoHandler := NewRepoHandler(
		services.DatabaseService,
		services.RenderService,
	)
	ingestHandler := NewIngestHandler(
		services.GitHubOIDCService,
		services.DatabaseService,
		services.RenderService,
	)

	return &Handlers{
		RepoHandler:    repoHandler,
		ImageHandler:   imageHandler,
		BadgeHandler:   badgeHandler,
		IngestHandler:  ingestHandler,
		HealthHandler:  healthHandler,
		DeleteHandler:  deleteHandler,
		SitemapHandler: sitemapHandler,
	}
}
