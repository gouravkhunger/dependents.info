package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/service"
	"dependents.info/pkg/utils"
)

type ImageHandler struct {
	dependentsService service.DependentsTasker
	store             service.Store
}

func NewImageHandler(
	store service.Store,
	dependentsService service.DependentsTasker,
) *ImageHandler {
	return &ImageHandler{
		store:             store,
		dependentsService: dependentsService,
	}
}

func (h *ImageHandler) SVGImage(c *fiber.Ctx) error {
	id := c.Query("id")
	repo := c.Params("owner") + "/" + c.Params("repo")

	name := repo
	if id != "" {
		name += ":" + id
	}

	var svg string
	err := h.store.Get("svg:"+name, &svg)

	if err != nil {
		taskErr := h.dependentsService.NewTask(c.UserContext(), repo, id, "image", func(total int, svg []byte) {
			_ = h.store.SaveWithTTL("svg:"+name, svg, 7*24*time.Hour)
			_ = h.store.SaveWithTTL("total:"+name, []byte(strconv.Itoa(total)), 7*24*time.Hour)
		})
		if taskErr != nil {
			return utils.SendError(c, fiber.StatusNotFound, "SVG image not found", taskErr)
		}
		err = h.store.Get("svg:"+name, &svg)
		if err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "SVG image not found", err)
		}
	}

	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")
	return c.Status(fiber.StatusOK).Type("svg").SendString(svg)
}
