package resources

import (
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS

//go:embed static/*
var StaticFiles embed.FS

func Templates() *template.Template {
	_templates, err := template.ParseFS(templateFiles, "templates/*")
	if err != nil {
		panic(fmt.Sprintf("Cannot parse templates: %v", err))
	}
	return _templates
}
