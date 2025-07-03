package handlers

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service/database"
	"dependents.info/internal/service/render"
	"dependents.info/pkg/utils"
)

type SitemapHandler struct {
	renderService   *render.RenderService
	databaseService *database.BadgerService
}

func NewSitemapHandler(
	dbService *database.BadgerService,
	renderService *render.RenderService,
) *SitemapHandler {
	return &SitemapHandler{
		databaseService: dbService,
		renderService:   renderService,
	}
}

func (h *SitemapHandler) Sitemap(c *fiber.Ctx) error {
	urls := make([]string, 0)
	h.databaseService.IterateKeys(func(key string) {
		urls = append(urls, key)
	})

	sitemapBytes, err := h.renderService.RenderSitemap(urls)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to render sitemap", err)
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("svg").Send(sitemapBytes)
}
