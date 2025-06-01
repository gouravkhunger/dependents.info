package image

import (
	"bytes"
	"fmt"

	svg "github.com/ajstarks/svgo"

	"dependents-img/internal/models"
)

func (i *ImageService) RenderSVG(data models.IngestRequest) []byte {
	w := bytes.NewBuffer(nil)
	canvas := svg.New(w)
	width, height := 400, 100
	canvas.Start(width, height)

	radius := 30
	x, y := 32, 32

	for i, dependent := range data.Dependents {
		offset := i * 48
		fillId := fmt.Sprintf("fill%d", i)

		canvas.Def()
		fmt.Fprintf(canvas.Writer, `<pattern id="%s" x="0%%" y="0%%" width="100%%" height="100%%" viewBox="0 0 512 512">`+"\n", fillId)
		fmt.Fprintf(canvas.Writer, `<image x="0%%" y="0%%" width="512" height="512" xlink:href="%s" />`+"\n", dependent.Image)
		canvas.PatternEnd()
		canvas.DefEnd()

		canvas.Circle(x+offset, y, radius, `stroke="#656c7688"`, `stroke-width="2"`, fmt.Sprintf(`fill="url(#%s)"`, fillId))
	}

	canvas.End()
	return w.Bytes()
}
