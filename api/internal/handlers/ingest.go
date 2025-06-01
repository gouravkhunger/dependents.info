package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"dependents-img/internal/config"
	"dependents-img/internal/env"
	"dependents-img/internal/models"
	"dependents-img/internal/service/database"
	"dependents-img/internal/service/github"
	"dependents-img/pkg/utils"
)

type IngestHandler struct {
	githubOIDCService *github.OIDCService
	databaseService   *database.BadgerService
}

func NewIngestHandler(githubOIDC *github.OIDCService, dbService *database.BadgerService) *IngestHandler {
	return &IngestHandler{
		databaseService:   dbService,
		githubOIDCService: githubOIDC,
	}
}

func (h *IngestHandler) Ingest(c *fiber.Ctx) error {
	config := config.FromContext(c.UserContext())
	name := fmt.Sprintf("%s/%s", c.Params("owner"), c.Params("repo"))

	if config == nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Configuration not found in context", nil)
	}

	if config.Environment == env.EnvProduction {
		token, err := utils.ExtractBearerToken(c.Get("Authorization"))
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "Invalid Authorization header", err)
		}

		if err := h.githubOIDCService.VerifyToken(c.Context(), token, name); err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "Repository ownership verification failed", err)
		}
	}

	var req models.IngestRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid JSON payload", err)
	}

	if err := req.Validate(); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid JSON payload", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to marshal request body", err)
	}

	if err = h.databaseService.Save(name, body); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to store data", err)
	}

	return utils.SendResponse(c, fiber.StatusOK, models.APIResponse{
		Success: true,
		Message: "Dependents data ingested successfully",
	})
}
