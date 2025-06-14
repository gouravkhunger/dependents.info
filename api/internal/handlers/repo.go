package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/models"
	"dependents.info/internal/service/database"
	"dependents.info/internal/service/render"
	"dependents.info/pkg/utils"
)

type RepoHandler struct {
	renderService   *render.RenderService
	databaseService *database.BadgerService
}

func NewRepoHandler(databaseService *database.BadgerService, renderService *render.RenderService) *RepoHandler {
	return &RepoHandler{
		renderService:   renderService,
		databaseService: databaseService,
	}
}

func (h *RepoHandler) RepoPage(c *fiber.Ctx) error {
	owner := c.Params("owner")
	repo := c.Params("repo")

	var total string
	err := h.databaseService.Get("total:"+owner+"/"+repo, &total)

	if err != nil {
		c.Set("X-Robots-Tag", "noindex, nofollow")
		return c.Redirect("https://github.com/"+owner+"/"+repo+"/network/dependents", fiber.StatusTemporaryRedirect)
	}

	totalInt, _ := strconv.Atoi(total)
	data := models.RepoPage{
		Total: totalInt,
		Owner: owner,
		Repo:  repo,
	}

	page, err := h.renderService.RenderPage(data)

	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Failed to generate repository page", err)
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("html").Send(page)
}
