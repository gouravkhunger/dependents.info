package github

import (
	"errors"
	"fmt"
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
	svc.fetchPage = func(url string) (string, error) {
		return dependentsPageHTML("owner/repo"), nil
	}

	var gotTotal int
	err := svc.NewTask("owner/repo", "", "badge", func(total int, svg []byte) {
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
	svc.fetchPage = func(url string) (string, error) {
		return dependentsPageHTML("owner/repo"), nil
	}
	svc.fetchImage = func(url string) (string, error) {
		return "data:image/png;base64,abc=", nil
	}

	var gotTotal int
	var gotSVG []byte
	err := svc.NewTask("owner/repo", "", "image", func(total int, svg []byte) {
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
	svc.fetchPage = func(url string) (string, error) {
		return "", errors.New("network error")
	}

	err := svc.NewTask("owner/repo", "", "badge", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewTask_RenderError(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{err: errors.New("render failed")})
	svc.fetchPage = func(url string) (string, error) {
		return dependentsPageHTML("owner/repo"), nil
	}
	svc.fetchImage = func(url string) (string, error) {
		return "data:image/png;base64,abc=", nil
	}

	err := svc.NewTask("owner/repo", "", "image", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewTask_WithPackageId(t *testing.T) {
	var capturedURL string
	svc := NewDependentsService(&stubRenderer{})
	svc.fetchPage = func(url string) (string, error) {
		capturedURL = url
		return dependentsPageHTML("owner/repo"), nil
	}

	if err := svc.NewTask("owner/repo", "pkg123", "badge", func(int, []byte) {}); err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}
	expected := "https://github.com/owner/repo/network/dependents?package_id=pkg123"
	if capturedURL != expected {
		t.Errorf("expected URL %q, got %q", expected, capturedURL)
	}
}

func TestNewTask_ImageFetchError_SkipsDependent(t *testing.T) {
	svc := NewDependentsService(&stubRenderer{result: []byte("<svg/>")})
	svc.fetchPage = func(url string) (string, error) {
		return dependentsPageHTML("owner/repo"), nil
	}
	svc.fetchImage = func(url string) (string, error) {
		return "", errors.New("image fetch failed")
	}

	var gotSVG []byte
	err := svc.NewTask("owner/repo", "", "image", func(total int, svg []byte) {
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
	callCount := 0
	svc.fetchImage = func(url string) (string, error) {
		callCount++
		return fmt.Sprintf("data:image/png;base64,%d", callCount), nil
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

	deps, err := svc.buildDependents(html)
	if err != nil {
		t.Fatalf("buildDependents() error = %v", err)
	}
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependents, got %d", len(deps))
	}
	if deps[0].Owner != "user2" || deps[1].Owner != "user1" || deps[2].Owner != "user3" {
		t.Errorf("expected star order user2, user1, user3; got %s, %s, %s", deps[0].Owner, deps[1].Owner, deps[2].Owner)
	}
	if callCount != 3 {
		t.Errorf("expected 3 image fetches, got %d", callCount)
	}
}
