package handlers

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service/database"
	"dependents.info/pkg/utils"
)

type ImageHandler struct {
	databaseService *database.BadgerService
}

func NewImageHandler(databaseService *database.BadgerService) *ImageHandler {
	return &ImageHandler{
		databaseService: databaseService,
	}
}

func (h *ImageHandler) SVGImage(c *fiber.Ctx) error {
	id := c.Query("id")
	name := c.Params("owner") + "/" + c.Params("repo")

	if id != "" {
		name += ":" + id
	}

	var svg string
	err := h.databaseService.Get("svg:"+name, &svg)

	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "SVG image not found", err)
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("svg").SendString(svg)
}
