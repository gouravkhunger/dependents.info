package handlers

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
	"dependents.info/internal/service"
	"dependents.info/pkg/utils"
)

type SitemapHandler struct {
	renderer service.Renderer
	store    service.Store
}

func NewSitemapHandler(
	store service.Store,
	renderer service.Renderer,
) *SitemapHandler {
	return &SitemapHandler{
		store:    store,
		renderer: renderer,
	}
}

func (h *SitemapHandler) Sitemap(c *fiber.Ctx) error {
	host := config.FromContext(c.UserContext()).Host()

	urls := make([]string, 0)
	seen := make(map[string]struct{})
	h.store.IterateKeys(func(key string) {
		route := utils.ToRoute(key)
		if _, exists := seen[route]; !exists {
			seen[route] = struct{}{}
			urls = append(urls, host+route)
		}
	})

	sitemapBytes, err := h.renderer.RenderSitemap(urls)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to render sitemap", err)
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("xml").Send(sitemapBytes)
}
