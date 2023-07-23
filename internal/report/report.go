//go:generate sh -c "yarn && yarn vite build"
package report

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"

	"github.com/alexbakker/gotchet/internal/format"
)

//go:embed dist/index.html
var tmplStr string
var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.New("report").Delims("#[", "]#").Parse(tmplStr)
	if err != nil {
		panic(err)
	}
}

func Render(c *format.TestCapture, w io.Writer) error {
	var buf bytes.Buffer
	be := base64.NewEncoder(base64.StdEncoding, &buf)
	zw := gzip.NewWriter(be)
	je := json.NewEncoder(zw)
	if err := je.Encode(c); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	if err := be.Close(); err != nil {
		return err
	}
	return tmpl.Execute(w, struct {
		Title string
		Data  template.HTML
	}{
		Title: "Go Test Report",
		Data:  template.HTML(buf.String()),
	})
}
