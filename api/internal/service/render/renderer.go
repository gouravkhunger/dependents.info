package render

import (
	"bytes"
	"embed"
	"fmt"
	html "html/template"
	text "text/template"

	"dependents.info/internal/models"
	"dependents.info/pkg/utils"
)

//go:embed templates/image.svg
var svgTemplate embed.FS

//go:embed templates/repo.html
var repoTemplate embed.FS

//go:embed templates/sitemap.xml
var sitemapTemplate embed.FS

var tmpl *text.Template
var repoTmpl *html.Template
var sitemapTmpl *text.Template
var funcMap *text.FuncMap

func init() {
	funcMap = &text.FuncMap{
		"formatNumber": utils.FormatNumber,
	}

	var err error
	tmpl, err = text.New("svg").Funcs(*funcMap).ParseFS(svgTemplate, "templates/image.svg")
	if err != nil {
		panic(fmt.Sprintf("failed to parse SVG template: %v", err))
	}

	sitemapTmpl, err = text.New("sitemap").ParseFS(sitemapTemplate, "templates/sitemap.xml")
	if err != nil {
		panic(fmt.Sprintf("failed to parse sitemap template: %v", err))
	}

	repoTmpl, err = html.New("repo").Funcs(*funcMap).ParseFS(repoTemplate, "templates/repo.html")
	if err != nil {
		panic(fmt.Sprintf("failed to parse repository template: %v", err))
	}
}

func (i *RenderService) RenderSVG(d models.IngestRequest) ([]byte, error) {
	w := bytes.NewBuffer(nil)

	data := &models.IngestRequest{
		Dependents: d.Dependents,
		Total:      d.Total - len(d.Dependents),
	}

	err := tmpl.ExecuteTemplate(w, "svg", data)
	return w.Bytes(), err
}

func (i *RenderService) RenderPage(d models.RepoPage) ([]byte, error) {
	w := bytes.NewBuffer(nil)

	err := repoTmpl.ExecuteTemplate(w, "repo", d)
	return w.Bytes(), err
}

func (i *RenderService) RenderSitemap(data []string) ([]byte, error) {
	w := bytes.NewBuffer(nil)

	err := sitemapTmpl.ExecuteTemplate(w, "sitemap", data)
	return w.Bytes(), err
}
