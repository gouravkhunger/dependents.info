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
	deleteHandler := NewDeleteHandler(services.Store)
	imageHandler := NewImageHandler(services.Store, services.DependentsService)
	badgeHandler := NewBadgeHandler(
		services.Store,
		services.DependentsService,
	)
	sitemapHandler := NewSitemapHandler(
		services.Store,
		services.Renderer,
	)
	repoHandler := NewRepoHandler(
		services.Store,
		services.Renderer,
	)
	ingestHandler := NewIngestHandler(
		services.OIDCVerifier,
		services.Store,
		services.Renderer,
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
