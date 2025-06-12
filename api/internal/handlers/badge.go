package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/models"
	"dependents.info/internal/service/database"
	"dependents.info/pkg/utils"
)

type BadgeHandler struct {
	databaseService *database.BadgerService
}

func NewBadgeHandler(databaseService *database.BadgerService) *BadgeHandler {
	return &BadgeHandler{
		databaseService: databaseService,
	}
}

func (h *BadgeHandler) Badge(c *fiber.Ctx) error {
	var total string
	name := c.Params("owner") + "/" + c.Params("repo")
	err := h.databaseService.Get("total:"+name, &total)

	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Total dependents not found", err)
	}

	// req := fiber.Get("https://img.shields.io/badge/users-" + total + "-" + h.color(total) + ".svg")

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return utils.SendResponse(c, fiber.StatusOK, models.APIResponse{
		Success: true,
		Message: "Service is healthy!",
	})
}

func (h *BadgeHandler) color(v string) string {
	total, _ := strconv.Atoi(v)
	switch {
	case total <= 0:
		return "red"
	case total < 10:
		return "yellow"
	case total < 100:
		return "yellowgreen"
	case total < 1000:
		return "green"
	default:
		return "brightgreen"
	}
}
