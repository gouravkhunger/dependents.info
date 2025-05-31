package middleware

import (
	"context"
	"dependents-img/internal/config"

	"github.com/gofiber/fiber/v2"
)

func Config(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), config.ConfigContextKey, cfg)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
