package handlers

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
	"dependents.info/internal/env"
	"dependents.info/internal/models"
	"dependents.info/internal/service/database"
	"dependents.info/internal/service/github"
	"dependents.info/internal/service/image"
	"dependents.info/pkg/utils"
)

type IngestHandler struct {
	githubOIDCService *github.OIDCService
	imageService      *image.ImageService
	databaseService   *database.BadgerService
}

func NewIngestHandler(
	githubOIDC *github.OIDCService,
	imageService *image.ImageService,
	dbService *database.BadgerService,
) *IngestHandler {
	return &IngestHandler{
		databaseService:   dbService,
		githubOIDCService: githubOIDC,
		imageService:      imageService,
	}
}

func (h *IngestHandler) Ingest(c *fiber.Ctx) error {
	config := config.FromContext(c.UserContext())
	name := c.Params("owner") + "/" + c.Params("repo")

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

	svgBytes, err := h.imageService.RenderSVG(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to render SVG", err)
	}

	if err := h.databaseService.Save("svg:"+name, svgBytes); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to store SVG", err)
	}

	return utils.SendResponse(c, fiber.StatusOK, models.APIResponse{
		Success: true,
		Message: "Dependents data ingested successfully",
	})
}
