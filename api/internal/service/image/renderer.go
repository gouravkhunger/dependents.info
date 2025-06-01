package image

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"dependents-img/internal/models"
)

//go:embed templates/image.svg
var svgTemplate embed.FS
var tmpl *template.Template

func init() {
	var err error
	//print cwd
	
	tmpl, err = template.New("svg").ParseFS(svgTemplate, "templates/image.svg")
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
