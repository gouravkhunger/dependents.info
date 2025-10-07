package github

import (
	"fmt"
	"io"
	"net/http"

	"dependents.info/internal/models"
	"dependents.info/internal/service/render"
	"dependents.info/pkg/utils"
)

type DependentsService struct {
	renderService *render.RenderService
}

func NewDependentsService(renderService *render.RenderService) *DependentsService {
	return &DependentsService{
		renderService: renderService,
	}
}

func (s *DependentsService) NewTask(repo string, id string, kind string, callback func(total int, svg []byte)) {
	url := "https://github.com/" + repo + "/network/dependents"
	if id != "" {
		url += "?package_id=" + id
	}
	page, err := fetchPage(url)
	if err != nil {
		return
	}
	total, err := utils.ParseTotalDependents(page, repo)
	if err != nil {
		return
	}
	if kind == "badge" {
		if callback != nil {
			callback(total, nil)
		}
		return
	}
	dependents, err := utils.ParseDependents(page)
	if err != nil {
		return
	}
	req := models.IngestRequest{
		Dependents: dependents,
		Total:      total,
	}
	svgBytes, err := s.renderService.RenderSVG(req)
	if err != nil {
		return
	}
	if callback != nil {
		callback(total, svgBytes)
	}
}

func fetchPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch dependents: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	return string(body), nil
}
