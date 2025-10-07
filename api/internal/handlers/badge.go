package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service/database"
	"dependents.info/internal/service/github"
	"dependents.info/pkg/utils"
)

type BadgeHandler struct {
	dependentsService *github.DependentsService
	databaseService   *database.BadgerService
}

func NewBadgeHandler(
	databaseService *database.BadgerService,
	dependentsService *github.DependentsService,
) *BadgeHandler {
	return &BadgeHandler{
		databaseService:   databaseService,
		dependentsService: dependentsService,
	}
}

func (h *BadgeHandler) Badge(c *fiber.Ctx) error {
	id := c.Query("id")
	repo := c.Params("repo")
	owner := c.Params("owner")

	name := owner + "/" + repo
	if id != "" {
		name += ":" + id
	}

	var total string
	err := h.databaseService.Get("total:"+name, &total)

	if err != nil {
		h.dependentsService.NewTask(owner, repo, id, "badge", func(total int, svg []byte) {
			h.databaseService.SaveWithTTL("total:"+name, []byte(strconv.Itoa(total)), 7*24*time.Hour)
		})
		err = h.databaseService.Get("total:"+name, &total)
		if err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "Total dependents not found", err)
		}
	}

	totalInt, _ := strconv.Atoi(total)
	u := "https://img.shields.io/badge/dependents-" + utils.FormatNumber(totalInt) + "-" + color(total)
	url := utils.SetParams(u, map[string]string{
		"logo":       c.Query("logo"),
		"label":      c.Query("label"),
		"style":      c.Query("style"),
		"color":      c.Query("color"),
		"logoColor":  c.Query("logoColor"),
		"labelColor": c.Query("labelColor"),
	})

	statusCode, body, errs := fiber.Get(url).Bytes()

	if len(errs) > 0 {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch badge image", errs[0])
	}

	if statusCode != fiber.StatusOK {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch badge image", nil)
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("svg").Send(body)
}

func color(v string) string {
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
