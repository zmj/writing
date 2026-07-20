package main

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*.html
var tmplFS embed.FS

//go:embed static/style.css static/CNAME
var staticFS embed.FS

var tmpl = template.Must(template.ParseFS(tmplFS, "templates/*.html"))

type indexData struct {
	Months []Month
}

func renderPage(name string, data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
