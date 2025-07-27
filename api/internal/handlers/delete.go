package handlers

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
	"dependents.info/internal/service/database"
	"dependents.info/pkg/utils"
)

type DeleteHandler struct {
	databaseService *database.BadgerService
}

func NewDeleteHandler(databaseService *database.BadgerService) *DeleteHandler {
	return &DeleteHandler{
		databaseService: databaseService,
	}
}

func (h *DeleteHandler) Delete(c *fiber.Ctx) error {
	cfg := config.FromContext(c.UserContext())

	token, err := utils.ExtractBearerToken(c.Get("Authorization"))
	if err != nil || token != cfg.Password {
		return utils.SendError(c, fiber.StatusUnauthorized, "Invalid Authorization header", err)
	}

	id := c.Query("id")
	name := c.Params("owner") + "/" + c.Params("repo")

	if id != "" {
		name += ":" + id
	}

	keys := []string{"total:" + name, "svg:" + name}
	for _, key := range keys {
		err = h.databaseService.Delete(key)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Failed to delete "+key, err)
		}
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
