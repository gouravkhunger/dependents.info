package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service"
	"dependents.info/pkg/utils"
)

type BadgeHandler struct {
	dependentsService service.DependentsTasker
	store             service.Store
}

func NewBadgeHandler(
	store service.Store,
	dependentsService service.DependentsTasker,
) *BadgeHandler {
	return &BadgeHandler{
		store:             store,
		dependentsService: dependentsService,
	}
}

func (h *BadgeHandler) resolveTotal(repo, id string) (string, error) {
	name := repo
	if id != "" {
		name += ":" + id
	}

	var total string
	err := h.store.Get("total:"+name, &total)
	if err != nil {
		taskErr := h.dependentsService.NewTask(repo, id, "badge", func(total int, svg []byte) {
			_ = h.store.SaveWithTTL("total:"+name, []byte(strconv.Itoa(total)), 7*24*time.Hour)
		})
		if taskErr != nil {
			return "", taskErr
		}
		err = h.store.Get("total:"+name, &total)
	}
	return total, err
}

func (h *BadgeHandler) Badge(c *fiber.Ctx) error {
	id := c.Query("id")
	repo := c.Params("owner") + "/" + c.Params("repo")

	total, err := h.resolveTotal(repo, id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Total dependents not found", err)
	}
	body, err := getBadge(total, c.Queries())
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to fetch badge image", err)
	}
	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("svg").Send(body)
}

func (h *BadgeHandler) Shields(c *fiber.Ctx) error {
	id := c.Query("id")
	repo := c.Params("owner") + "/" + c.Params("repo")

	total, err := h.resolveTotal(repo, id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Total dependents not found", err)
	}

	totalInt, _ := strconv.Atoi(total)
	label := c.Query("label", "dependents")
	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"schemaVersion": 1,
		"label":         label,
		"message":       utils.FormatNumber(totalInt),
		"color":         color(total),
		"cacheSeconds":  86400,
	})
}

func (h *BadgeHandler) SelfBadge(c *fiber.Ctx) error {
	var total string
	seen := make(map[string]struct{})
	h.store.IterateKeys(func(key string) {
		route := utils.ToRoute(key)
		if _, exists := seen[route]; !exists {
			seen[route] = struct{}{}
		}
	})
	total = strconv.Itoa(len(seen))
	queries := c.Queries()
	queries["label"] = "users"
	body, err := getBadge(total, queries)
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
