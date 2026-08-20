package github

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"dependents.info/internal/models"
	"dependents.info/pkg/utils"
)

type renderer interface {
	RenderSVG(d models.IngestRequest) ([]byte, error)
}

type DependentsService struct {
	renderer   renderer
	client     *http.Client
	fetchPage  func(url string) (string, error)
	fetchImage func(url string) (string, error)
}

func NewDependentsService(r renderer) *DependentsService {
	s := &DependentsService{
		renderer: r,
		client:   &http.Client{Timeout: 8 * time.Second},
	}
	s.fetchPage = s.defaultFetchPage
	s.fetchImage = s.defaultFetchImage
	return s
}

func (s *DependentsService) NewTask(repo string, id string, kind string, callback func(total int, svg []byte)) error {
	url := "https://github.com/" + repo + "/network/dependents"
	if id != "" {
		url += "?package_id=" + id
	}
	page, err := s.fetchPage(url)
	if err != nil {
		return fmt.Errorf("fetch page: %w", err)
	}
	total, err := utils.ParseTotalDependents(page, repo)
	if err != nil {
		return fmt.Errorf("parse total: %w", err)
	}
	if kind == "badge" {
		if callback != nil {
			callback(total, nil)
		}
		return nil
	}
	dependents, err := s.buildDependents(page)
	if err != nil {
		return fmt.Errorf("build dependents: %w", err)
	}
	req := models.IngestRequest{
		Dependents: dependents,
		Total:      total,
	}
	svgBytes, err := s.renderer.RenderSVG(req)
	if err != nil {
		return fmt.Errorf("render svg: %w", err)
	}
	if callback != nil {
		callback(total, svgBytes)
	}
	return nil
}

func (s *DependentsService) buildDependents(page string) ([]models.Dependent, error) {
	nodes, err := utils.ParseDependentNodes(page)
	if err != nil {
		return nil, err
	}

	type result struct {
		index int
		dep   models.Dependent
		err   error
	}

	var wg sync.WaitGroup
	results := make(chan result, len(nodes))

	for i, node := range nodes {
		wg.Add(1)
		go func(idx int, n utils.DependentInfo) {
			defer wg.Done()
			image, err := s.fetchImage(n.ImageURL)
			results <- result{
				index: idx,
				dep:   models.Dependent{Image: image, Stars: n.Stars, Owner: n.Owner},
				err:   err,
			}
		}(i, node)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	byIndex := make([]models.Dependent, len(nodes))
	ok := make([]bool, len(nodes))
	for r := range results {
		if r.err == nil {
			byIndex[r.index] = r.dep
			ok[r.index] = true
		}
	}

	dependents := make([]models.Dependent, 0, len(nodes))
	for i, d := range byIndex {
		if ok[i] {
			dependents = append(dependents, d)
		}
	}
	return dependents, nil
}

func (s *DependentsService) defaultFetchPage(url string) (string, error) {
	var body string
	err := utils.RetryWithBackoff(3, 500*time.Millisecond, func() error {
		resp, err := s.client.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch page: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 500 {
			return fmt.Errorf("server error: status %d", resp.StatusCode)
		}
		if resp.StatusCode != http.StatusOK {
			io.Copy(io.Discard, resp.Body)
			return utils.PermanentError{Err: fmt.Errorf("unexpected status %d", resp.StatusCode)}
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}
		body = string(data)
		return nil
	})
	return body, err
}

func (s *DependentsService) defaultFetchImage(url string) (string, error) {
	var dataURI string
	err := utils.RetryWithBackoff(3, 300*time.Millisecond, func() error {
		resp, err := s.client.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch image: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 500 {
			return fmt.Errorf("server error: status %d", resp.StatusCode)
		}
		if resp.StatusCode != http.StatusOK {
			io.Copy(io.Discard, resp.Body)
			return utils.PermanentError{Err: fmt.Errorf("unexpected status %d", resp.StatusCode)}
		}
		imageData, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read image data: %w", err)
		}
		base64Str := base64.StdEncoding.EncodeToString(imageData)
		mimeType := resp.Header.Get("Content-Type")
		dataURI = fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
		return nil
	})
	return dataURI, err
}
