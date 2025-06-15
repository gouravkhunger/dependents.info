package utils

import (
	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/models"
)

func SendResponse(c *fiber.Ctx, status int, response models.APIResponse) error {
	return c.Status(status).JSON(response)
}

func SendError(c *fiber.Ctx, status int, message string, err error) error {
	response := models.APIResponse{
		Success: false,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	c.Set(fiber.HeaderXRobotsTag, "noindex, nofollow")
	return c.Status(status).JSON(response)
}
