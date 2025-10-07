package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service/database"
	"dependents.info/internal/service/github"
	"dependents.info/pkg/utils"
)

type ImageHandler struct {
	dependentsService *github.DependentsService
	databaseService   *database.BadgerService
}

func NewImageHandler(
	databaseService *database.BadgerService,
	dependentsService *github.DependentsService,
) *ImageHandler {
	return &ImageHandler{
		databaseService:   databaseService,
		dependentsService: dependentsService,
	}
}

func (h *ImageHandler) SVGImage(c *fiber.Ctx) error {
	id := c.Query("id")
	repo := c.Params("repo")
	owner := c.Params("owner")

	name := owner + "/" + repo
	if id != "" {
		name += ":" + id
	}

	var svg string
	err := h.databaseService.Get("svg:"+name, &svg)

	if err != nil {
		h.dependentsService.NewTask(owner, repo, id, "image", func(total int, svg []byte) {
			h.databaseService.SaveWithTTL("svg:"+name, svg, 7*24*time.Hour)
			h.databaseService.SaveWithTTL("total:"+name, []byte(strconv.Itoa(total)), 7*24*time.Hour)
		})
		err = h.databaseService.Get("svg:"+name, &svg)
		if err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "SVG image not found", err)
		}
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("svg").SendString(svg)
}
