package utils

import (
	"embed"
	"fmt"
	"math"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var units = []struct {
	value  int
	suffix string
}{
	{1_000_000_000, "B"},
	{1_000_000, "M"},
	{1_000, "K"},
}

func round(x float64) float64 {
	return math.Floor(x + 0.5)
}

func digits(x float64) int {
	return int(math.Log10(x)) + 1
}

func limitFn(x int) int {
	digits := int(math.Log10(float64(x))) + 1
	return int(math.Pow(10, float64(digits+1))) - 1
}

func ValidateRepository(value string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9-]+/[a-zA-Z0-9._-]+$`)
	return re.MatchString(value)
}

func FormatNumber(value int) string {
	for i, u := range units {
		if value >= u.value {
			limit := limitFn(u.value)
			div := float64(value) / float64(u.value)
			var res float64
			if value < limit {
				res = round(div*10) / 10
			} else {
				res = round(div)
			}
			if digits(res) > 3 {
				next := units[max(0, i-1)]
				return fmt.Sprintf("%.0f%s", round(res/1_000), next.suffix)
			}
			if res == math.Floor(res) {
				return fmt.Sprintf("%.0f%s", res, u.suffix)
			}
			return fmt.Sprintf("%.1f%s", res, u.suffix)
		}
	}

	return fmt.Sprintf("%d", value)
}

func LoadStylesFile(static *embed.FS) {
	files, _ := static.ReadDir("static/assets")
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".css") {
			os.Setenv("STYLES_FILE", file.Name())
			break
		}
	}
	if stylesFile := os.Getenv("STYLES_FILE"); stylesFile == "" {
		panic("No CSS file found in static/assets directory")
	}
}

func ToRoute(key string) string {
	parts := strings.FieldsFunc(key, func(r rune) bool {
		return r == ':'
	})
	if len(parts) < 2 || !strings.Contains(parts[1], "/") {
		return ""
	}
	url := "/" + parts[1]
	if len(parts) > 2 {
		url += fmt.Sprintf("?id=%s", parts[2])
	}
	return url
}

func SetParams(s string, params map[string]string) string {
	url, _ := url.Parse(s)
	q := url.Query()
	for key, value := range params {
		if value != "" {
			q.Set(key, value)
		}
	}
	url.RawQuery = q.Encode()
	return url.String()
}

func ExtractBearerToken(authHeader string) (string, error) {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return strings.TrimSpace(authHeader[len(prefix):]), nil
	}
	return "", fiber.ErrUnauthorized
}
