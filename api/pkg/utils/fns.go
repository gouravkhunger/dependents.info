package utils

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"dependents.info/internal/models"
	"github.com/andybalholm/cascadia"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/html"
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

func ParseTotalDependents(doc string, repo string) (int, error) {
	node, err := html.Parse(strings.NewReader(doc))
	if err != nil {
		return 0, fmt.Errorf("failed to parse HTML: %w", err)
	}
	selector, err := cascadia.Compile(
		fmt.Sprintf(`a[href*="/%s/network/dependents?dependent_type=REPOSITORY"]`, repo),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to compile selector: %w", err)
	}
	anchor := cascadia.Query(node, selector)
	if anchor == nil {
		return 0, fmt.Errorf("could not find anchor for %s", repo)
	}
	number, err := extractNumber(anchor)
	if err != nil {
		return 0, fmt.Errorf("failed to extract number: %w", err)
	}
	return number, nil
}

func ParseDependents(doc string) ([]models.Dependent, error) {
	node, err := html.Parse(strings.NewReader(doc))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	imgSel, _ := cascadia.Compile("img")
	dependentSel, _ := cascadia.Compile(`[data-test-id="dg-repo-pkg-dependent"]`)
	starsSel, _ := cascadia.Compile(".octicon-star")
	nodes := cascadia.QueryAll(node, dependentSel)
	dependents := make([]models.Dependent, 0, 30)
	for _, el := range nodes {
		var image string
		imgNode := cascadia.Query(el, imgSel)
		starsNode := cascadia.Query(el, starsSel)
		if imgNode == nil {
			continue
		}
		image, err = imageNodeToUrl(imgNode)
		if err != nil {
			continue
		}
		image, err = imageUrlToBase64(image)
		if err != nil {
			continue
		}
		stars, _ := extractNumber(starsNode.Parent)
		dependents = append(dependents, models.Dependent{
			Image: image,
			Stars: stars,
		})
	}
	sort.Slice(dependents, func(i, j int) bool {
		return dependents[i].Stars > dependents[j].Stars
	})
	if len(dependents) > 11 {
		return dependents[:11], nil
	}
	return dependents, nil
}

func extractNumber(anchor *html.Node) (int, error) {
	var sb strings.Builder
	for c := anchor.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			sb.WriteString(c.Data)
		}
	}
	text := strings.TrimSpace(sb.String())
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return 0, fmt.Errorf("no text found in anchor")
	}
	numString := strings.ReplaceAll(fields[0], ",", "")
	number, _ := strconv.Atoi(numString)
	return number, nil
}

func imageNodeToUrl(n *html.Node) (string, error) {
	for _, attr := range n.Attr {
		if attr.Key == "src" {
			return SetParams(attr.Val, map[string]string{"s": "100"}), nil
		}
	}
	return "", fmt.Errorf("no src attribute found in image node")
}

func imageUrlToBase64(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching image: %w", err)
	}
	defer resp.Body.Close()
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading image data: %w", err)
	}
	base64Str := base64.StdEncoding.EncodeToString(imageData)
	mimeType := resp.Header.Get("Content-Type")
	dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
	return dataURI, nil
}
