package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/radoslawg/video_manager/resources"
)

var templates *template.Template

func main() {
	templates = resources.Templates()

	server := http.NewServeMux()
	server.HandleFunc("/", serveIndex)
	server.Handle("/static/", http.FileServer(http.FS(resources.StaticFiles)))

	http.ListenAndServe(":8080", server)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	err := templates.Execute(w, "index.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Cannot serve template %v", err))
	}
}
