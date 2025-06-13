package handlers

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/models"
	"dependents.info/internal/service/render"
	"dependents.info/pkg/utils"
)

type RepoHandler struct {
	renderService *render.RenderService
}

func NewRepoHandler(renderService *render.RenderService) *RepoHandler {
	return &RepoHandler{
		renderService: renderService,
	}
}

func (h *RepoHandler) RepoPage(c *fiber.Ctx) error {
	data := models.RepoPage{
		Owner: c.Params("owner"),
		Repo:  c.Params("repo"),
	}

	page, err := h.renderService.RenderPage(data)

	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Failed to generate repository page", err)
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("html").Send(page)
}
