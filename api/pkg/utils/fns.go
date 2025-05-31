package utils

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ValidateRepository(value string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9-]+/[a-zA-Z0-9._-]+$`)
	return re.MatchString(value)
}

func ExtractBearerToken(authHeader string) (string, error) {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return strings.TrimSpace(authHeader[len(prefix):]), nil
	}
	return "", fiber.ErrUnauthorized
}
