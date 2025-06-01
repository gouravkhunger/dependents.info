package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
)

func Config(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), config.ConfigContextKey, cfg)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
