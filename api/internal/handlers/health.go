package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"dependents-img/internal/models"
	"dependents-img/pkg/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	return utils.SendResponse(c, fiber.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
		Message: "Service is healthy!",
	})
}
