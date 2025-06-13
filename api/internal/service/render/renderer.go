package render

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"dependents.info/internal/models"
	"dependents.info/pkg/utils"
)

//go:embed templates/image.svg
var svgTemplate embed.FS

var tmpl *template.Template
var funcMap *template.FuncMap

func init() {
	funcMap = &template.FuncMap{
		"formatNumber": utils.FormatNumber,
	}

	var err error
	tmpl, err = template.New("svg").Funcs(*funcMap).ParseFS(svgTemplate, "templates/image.svg")
	if err != nil {
		panic(fmt.Sprintf("failed to parse SVG template: %v", err))
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
