package image

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"golang.org/x/text/message"

	"dependents-img/internal/models"
)

//go:embed templates/image.svg
var svgTemplate embed.FS

var p *message.Printer
var tmpl *template.Template
var funcMap *template.FuncMap

func init() {
	p = message.NewPrinter(message.MatchLanguage("en"))

	funcMap = &template.FuncMap{
		"formatNumber": func(n int) string {
			return p.Sprintf("%d", n)
		},
	}

	var err error
	tmpl, err = template.New("svg").Funcs(*funcMap).ParseFS(svgTemplate, "templates/image.svg")
	if err != nil {
		panic(fmt.Sprintf("failed to parse SVG template: %v", err))
	}
}

func (i *ImageService) RenderSVG(d models.IngestRequest) ([]byte, error) {
	w := bytes.NewBuffer(nil)

	data := &models.IngestRequest{
		Dependents: d.Dependents,
		Total:      d.Total - len(d.Dependents),
	}

	err := tmpl.ExecuteTemplate(w, "svg", data)
	return w.Bytes(), err
}
