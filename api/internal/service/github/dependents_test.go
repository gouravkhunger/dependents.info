package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"dependents.info/internal/models"
)

type stubRenderer struct {
	result []byte
	err    error
}

func (s *stubRenderer) RenderSVG(_ models.IngestRequest) ([]byte, error) {
	return s.result, s.err
}

func dependentsPageHTML(repo string) string {
	return fmt.Sprintf(`<html><body>
		<div role="status" class="table-list-header-toggle states flex-auto pl-0">
			<a href="/%s/network/dependents?dependent_type=REPOSITORY">
				100 Repositories
			</a>
		</div>
		<div data-test-id="dg-repo-pkg-dependent">
			<a data-hovercard-type="user">testowner</a>
			<img src="https://avatars.githubusercontent.com/u/1?v=4" />
			<span class="octicon-star"></span>
			42
		</div>
	</body></html>`, repo)
}

func TestNewTask_Badge(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{})
	svc.fetchPage = func(_ context.Context, url string) (string, error) {
		return dependentsPageHTML("owner/repo"), nil
	}

	var gotTotal int
	err := svc.NewTask(context.Background(), "owner/repo", "", "badge", func(total int, svg []byte) {
		gotTotal = total
	})
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}
	if gotTotal != 100 {
		t.Errorf("expected total 100, got %d", gotTotal)
	}
}

func TestNewTask_Image(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{result: []byte("<svg>rendered</svg>")})
	svc.fetchPage = func(_ context.Context, url string) (string, error) {
		return dependentsPageHTML("owner/repo"), nil
	}
	svc.fetchImage = func(_ context.Context, url string) (string, error) {
		return "data:image/png;base64,abc=", nil
	}

	var gotTotal int
	var gotSVG []byte
	err := svc.NewTask(context.Background(), "owner/repo", "", "image", func(total int, svg []byte) {
		gotTotal = total
		gotSVG = svg
	})
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}
	if gotTotal != 100 {
		t.Errorf("expected total 100, got %d", gotTotal)
	}
	if string(gotSVG) != "<svg>rendered</svg>" {
		t.Errorf("expected rendered SVG, got %q", string(gotSVG))
	}
}

func TestNewTask_FetchPageError(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{})
	svc.fetchPage = func(_ context.Context, url string) (string, error) {
		return "", errors.New("network error")
	}

	err := svc.NewTask(context.Background(), "owner/repo", "", "badge", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewTask_RenderError(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{err: errors.New("render failed")})
	svc.fetchPage = func(_ context.Context, url string) (string, error) {
		return dependentsPageHTML("owner/repo"), nil
	}
	svc.fetchImage = func(_ context.Context, url string) (string, error) {
		return "data:image/png;base64,abc=", nil
	}

	err := svc.NewTask(context.Background(), "owner/repo", "", "image", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewTask_WithPackageId(t *testing.T) {
	var capturedURL string
	svc := NewDependentsService(&stubRenderer{})
	svc.fetchPage = func(_ context.Context, url string) (string, error) {
		capturedURL = url
		return dependentsPageHTML("owner/repo"), nil
	}

	if err := svc.NewTask(context.Background(), "owner/repo", "pkg123", "badge", func(int, []byte) {}); err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}
	expected := "https://github.com/owner/repo/network/dependents?package_id=pkg123"
	if capturedURL != expected {
		t.Errorf("expected URL %q, got %q", expected, capturedURL)
	}
}

func TestNewTask_ImageFetchError_SkipsDependent(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{result: []byte("<svg/>")})
	svc.fetchPage = func(_ context.Context, url string) (string, error) {
		return dependentsPageHTML("owner/repo"), nil
	}
	svc.fetchImage = func(_ context.Context, url string) (string, error) {
		return "", errors.New("image fetch failed")
	}

	var gotSVG []byte
	err := svc.NewTask(context.Background(), "owner/repo", "", "image", func(total int, svg []byte) {
		gotSVG = svg
	})
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}
	if string(gotSVG) != "<svg/>" {
		t.Errorf("expected rendered SVG even with failed images, got %q", string(gotSVG))
	}
}

func TestBuildDependents_Parallel(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{})
	var callCount atomic.Int32
	svc.fetchImage = func(_ context.Context, url string) (string, error) {
		callCount.Add(1)
		return "data:image/png;base64," + url, nil
	}

	html := `<html><body>
		<div role="status"><a href="/o/r/network/dependents?dependent_type=REPOSITORY">10 Repositories</a></div>
		<div data-test-id="dg-repo-pkg-dependent">
			<a data-hovercard-type="user">user1</a><img src="https://example.com/1.png"/><span class="octicon-star"></span> 10
		</div>
		<div data-test-id="dg-repo-pkg-dependent">
			<a data-hovercard-type="user">user2</a><img src="https://example.com/2.png"/><span class="octicon-star"></span> 20
		</div>
		<div data-test-id="dg-repo-pkg-dependent">
			<a data-hovercard-type="user">user3</a><img src="https://example.com/3.png"/><span class="octicon-star"></span> 5
		</div>
	</body></html>`

	deps, err := svc.buildDependents(context.Background(), html)
	if err != nil {
		t.Fatalf("buildDependents() error = %v", err)
	}
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependents, got %d", len(deps))
	}
	if deps[0].Owner != "user2" || deps[1].Owner != "user1" || deps[2].Owner != "user3" {
		t.Errorf("expected star order user2, user1, user3; got %s, %s, %s", deps[0].Owner, deps[1].Owner, deps[2].Owner)
	}
	if !strings.Contains(deps[0].Image, "https://example.com/2.png") ||
		!strings.Contains(deps[1].Image, "https://example.com/1.png") ||
		!strings.Contains(deps[2].Image, "https://example.com/3.png") {
		t.Errorf("images not aligned with owners: %q %q %q", deps[0].Image, deps[1].Image, deps[2].Image)
	}
	if callCount.Load() != 3 {
		t.Errorf("expected 3 image fetches, got %d", callCount.Load())
	}
}

func TestBuildDependents_BackfillsFailedAvatars(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{})
	svc.fetchImage = func(_ context.Context, url string) (string, error) {
		if strings.Contains(url, "/11.png") {
			return "", errors.New("flaky")
		}
		return "data:image/png;base64," + url, nil
	}

	var html strings.Builder
	html.WriteString(`<html><body>`)
	for i := range 12 {
		fmt.Fprintf(&html, `<div data-test-id="dg-repo-pkg-dependent">
			<a data-hovercard-type="user">user%d</a>
			<img src="https://example.com/%d.png"/>
			<span class="octicon-star"></span> %d
		</div>`, i, i, i)
	}
	html.WriteString(`</body></html>`)

	deps, err := svc.buildDependents(context.Background(), html.String())
	if err != nil {
		t.Fatalf("buildDependents() error = %v", err)
	}
	if len(deps) != 11 {
		t.Fatalf("expected 11 dependents, got %d", len(deps))
	}
	if deps[0].Owner != "user10" {
		t.Errorf("expected highest remaining star owner user10, got %s", deps[0].Owner)
	}
	for _, d := range deps {
		if d.Owner == "user11" || strings.Contains(d.Image, "/11.png") {
			t.Fatalf("failed avatar should be skipped, got %+v", d)
		}
	}
}

func TestDefaultFetchPage_Retries429(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits.Add(1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		io.WriteString(w, "ok")
	}))
	t.Cleanup(srv.Close)

	svc := NewDependentsService(&stubRenderer{})
	body, err := svc.defaultFetchPage(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("defaultFetchPage() error = %v", err)
	}
	if body != "ok" {
		t.Errorf("body = %q", body)
	}
	if hits.Load() != 2 {
		t.Errorf("expected 2 hits, got %d", hits.Load())
	}
}
