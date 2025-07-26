package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
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
	id := c.Query("id")
	repo := c.Params("repo")
	owner := c.Params("owner")
	name := owner + "/" + repo
	cfg := config.FromContext(c.UserContext())

	if id != "" {
		name += ":" + id
	}

	var total string
	err := h.databaseService.Get("total:"+name, &total)

	if err != nil {
		url := "https://github.com/" + owner + "/" + repo + "/network/dependents"
		if id != "" {
			url += "?package_id=" + id
		}
		c.Set(fiber.HeaderXRobotsTag, "noindex, nofollow")
		return c.Redirect(url, fiber.StatusTemporaryRedirect)
	}

	var image string
	err = h.databaseService.Get("svg:"+name, &image)
	if err != nil {
		image = ""
	}

	totalInt, _ := strconv.Atoi(total)
	data := models.RepoPage{
		StylesFile: cfg.StylesFile,
		HasImage:   image != "",
		Total:      totalInt,
		Owner:      owner,
		Repo:       repo,
		Id:         id,
	}

	page, err := h.renderService.RenderPage(data)

	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Failed to generate repository page", err)
	}

	if id != "" {
		c.Set(fiber.HeaderXRobotsTag, "noindex, nofollow")
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("html").Send(page)
}
