package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"dependents.info/internal/config"
	"dependents.info/internal/models"
	"dependents.info/internal/service/database"
	"dependents.info/internal/service/render"
	"dependents.info/pkg/utils"
)

type RepoHandler struct {
	renderService   *render.RenderService
	databaseService *database.BadgerService
}

func NewRepoHandler(databaseService *database.BadgerService, renderService *render.RenderService) *RepoHandler {
	return &RepoHandler{
		renderService:   renderService,
		databaseService: databaseService,
	}
}

func splitRepoFormat(raw string) (repo, format string) {
	switch {
	case strings.HasSuffix(raw, ".json"):
		return strings.TrimSuffix(raw, ".json"), "json"
	case strings.HasSuffix(raw, ".md"):
		return strings.TrimSuffix(raw, ".md"), "md"
	default:
		return raw, "html"
	}
}

func (h *RepoHandler) RepoPage(c *fiber.Ctx) error {
	id := c.Query("id")
	owner := c.Params("owner")
	repo, format := splitRepoFormat(c.Params("repo"))
	name := owner + "/" + repo
	cfg := config.FromContext(c.UserContext())

	if id != "" {
		name += ":" + id
	}

	var total string
	err := h.databaseService.Get("total:"+name, &total)

	if err != nil {
		if format != "html" {
			return utils.SendError(c, fiber.StatusNotFound, "Total dependents not found", err)
		}
		url := "https://github.com/" + owner + "/" + repo + "/network/dependents"
		if id != "" {
			url += "?package_id=" + id
		}
		c.Set(fiber.HeaderXRobotsTag, "noindex, nofollow")
		return c.Redirect(url, fiber.StatusTemporaryRedirect)
	}

	var image string
	err = h.databaseService.Get("svg:"+name, &image)
	if err != nil {
		image = ""
	}

	totalInt, _ := strconv.Atoi(total)
	hasImage := image != ""
	c.Set(fiber.HeaderCacheControl, "public, max-age=86400, must-revalidate")

	if format == "json" {
		return c.Status(fiber.StatusOK).JSON(repoJSON(cfg.Host(), owner, repo, id, totalInt, hasImage))
	}
	if format == "md" {
		c.Type("md")
		return c.Status(fiber.StatusOK).SendString(repoMarkdown(cfg.Host(), owner, repo, id, totalInt, hasImage))
	}

	data := models.RepoPage{
		StylesFile: cfg.StylesFile,
		HasImage:   hasImage,
		Total:      totalInt,
		Owner:      owner,
		Repo:       repo,
		Id:         id,
	}

	page, err := h.renderService.RenderPage(data)

	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Failed to generate repository page", err)
	}

	if id != "" {
		c.Set(fiber.HeaderXRobotsTag, "noindex, nofollow")
	}

	return c.Status(fiber.StatusOK).Type("html").Send(page)
}

func repoJSON(host, owner, repo, id string, total int, hasImage bool) fiber.Map {
	q := ""
	if id != "" {
		q = "?id=" + id
	}
	name := owner + "/" + repo
	gh := "https://github.com/" + name + "/network/dependents"
	if id != "" {
		gh += "?package_id=" + id
	}
	payload := fiber.Map{
		"owner":     owner,
		"repo":      repo,
		"total":     total,
		"has_image": hasImage,
		"urls": fiber.Map{
			"page":              host + "/" + name + q,
			"badge":             host + "/" + name + "/badge" + q,
			"image":             host + "/" + name + "/image" + q,
			"markdown":          host + "/" + name + ".md" + q,
			"json":              host + "/" + name + ".json" + q,
			"shields":           host + "/" + name + "/shields.json" + q,
			"github_dependents": gh,
		},
	}
	if id != "" {
		payload["id"] = id
	}
	return payload
}

func repoMarkdown(host, owner, repo, id string, total int, hasImage bool) string {
	q := ""
	if id != "" {
		q = "?id=" + id
	}
	name := owner + "/" + repo
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", name)
	if id != "" {
		fmt.Fprintf(&b, "package `%s`\n\n", id)
	}
	fmt.Fprintf(&b, "**%s** GitHub network dependents.\n\n", utils.FormatNumber(total))
	fmt.Fprintf(&b, "[![dependents](%s/%s/badge%s)](%s/%s%s)\n", host, name, q, host, name, q)
	if hasImage {
		fmt.Fprintf(&b, "\n![used by](%s/%s/image%s)\n", host, name, q)
	}
	gh := "https://github.com/" + name + "/network/dependents"
	if id != "" {
		gh += "?package_id=" + id
	}
	fmt.Fprintf(&b, "\n- [dependents.info](%s/%s%s)\n", host, name, q)
	fmt.Fprintf(&b, "- [GitHub dependents](%s)\n", gh)
	fmt.Fprintf(&b, "- [JSON](%s/%s.json%s)\n", host, name, q)
	return b.String()
}
