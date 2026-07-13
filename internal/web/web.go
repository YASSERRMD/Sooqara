package web

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/* static/*
var fs embed.FS

// Templates holds parsed HTML templates.
var Templates *template.Template

// InitTemplates loads and parses all templates.
func InitTemplates() error {
	Templates = template.Must(template.New("").ParseFS(fs, "templates/*.html"))
	return nil
}

// ServeFiles serves static files from the embedded filesystem.
func ServeFiles(mux *http.ServeMux) {
	mux.Handle("GET /static/", http.FileServer(http.FS(fs)))
}
