package github

import (
	"fmt"
	"io"
	"net/http"

	"dependents.info/pkg/utils"
)

type DependentsService struct {}

func NewDependentsService() *DependentsService {
	return &DependentsService{}
}

func (s *DependentsService) GetTotalDependents(repo string, id string) (string, error) {
	url := "https://github.com/" + repo + "/network/dependents"
	if id != "" {
		url += "?package_id=" + id
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch dependents: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	return utils.ParseTotalDependents(string(body), repo)
}
