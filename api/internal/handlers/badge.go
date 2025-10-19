package handlers

import (
	"fmt"
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
	repo := c.Params("owner") + "/" + c.Params("repo")

	name := repo
	if id != "" {
		name += ":" + id
	}

	var total string
	err := h.databaseService.Get("total:"+name, &total)

	if err != nil {
		h.dependentsService.NewTask(repo, id, "badge", func(total int, svg []byte) {
			h.databaseService.SaveWithTTL("total:"+name, []byte(strconv.Itoa(total)), 7*24*time.Hour)
		})
		err = h.databaseService.Get("total:"+name, &total)
		if err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "Total dependents not found", err)
		}
	}
	body, err := getBadge(total, c.Queries())
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch badge image", err)
	}
	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("svg").Send(body)
}

func getBadge(total string, q map[string]string) ([]byte, error) {
	totalInt, _ := strconv.Atoi(total)
	u := "https://img.shields.io/badge/dependents-" + utils.FormatNumber(totalInt) + "-" + color(total)
	url := utils.SetParams(u, map[string]string{
		"logo":       q["logo"],
		"label":      q["label"],
		"style":      q["style"],
		"color":      q["color"],
		"logoColor":  q["logoColor"],
		"labelColor": q["labelColor"],
	})
	statusCode, body, errs := fiber.Get(url).Bytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	if statusCode != fiber.StatusOK {
		return nil, fmt.Errorf("failed to fetch badge: status code %d", statusCode)
	}
	return body, nil
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
