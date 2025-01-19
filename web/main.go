package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type FileLinks struct {
	FileName         string
	Links            []string
	OriginalFileName []string
	Titles           []string
}

const PATH = "/media/sda1/youtube/!nextdaily"

// const PATH = "z:/youtube/!nextdaily"

//go:embed templates/*
var content embed.FS
var templates *template.Template

func main() {
	_templates, err := template.ParseFS(content, "templates/*.gtpl")
	if err != nil {
		panic(err)
	}
	templates = _templates
	http.HandleFunc("/", listFilesHandler)
	http.HandleFunc("/view/", viewFileHandler)
	http.HandleFunc("/delete/", deleteLinkHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	dir := PATH // Directory containing the files
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mkv") {
			dates := strings.Split(strings.TrimSpace(string(file.Name())), "#")
			if !slices.Contains(fileNames, dates[0]) {
				fileNames = append(fileNames, dates[0])
			}
		}
	}

	// log.Println(templates.DefinedTemplates())
	templates.Lookup("index.gtpl").Execute(w, fileNames)
}

func viewFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/view/")
	fileName = filepath.Base(fileName) // Prevent directory traversal
	dir := PATH                        // Directory containing the files

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	var links []string
	var original_filenames []string
	var titles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mkv") && strings.HasPrefix(file.Name(), fileName) {
			i := 2
			ids := strings.Split(strings.TrimSpace(string(file.Name())), "#")
			if len(ids[1]) == 8 {
				_, err := strconv.ParseInt(ids[1], 10, 32)
				if err == nil {
					i = 3
				}
			}
			if !slices.Contains(links, ids[i]) {
				original_filenames = append(original_filenames, file.Name())
				links = append(links, ids[i])
				titles = append(titles, strings.ReplaceAll(ids[i-1], "_", " "))
			}
		}
	}

	if len(links) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	fileLinks := FileLinks{
		FileName:         fileName,
		Links:            links,
		OriginalFileName: original_filenames,
		Titles:           titles,
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	w.Header().Set("Expires", "0")                                         // Proxies.
	tmpl := templates.Lookup("day_view.gtpl")
	tmpl.Execute(w, fileLinks)
}

func deleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/delete/"), "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fileName, err := url.QueryUnescape(parts[0]) // Prevent directory traversal
	if err != nil {
		http.Error(w, "Unable to parse filename", http.StatusInternalServerError)
		return
	}
	view := parts[1]
	dir := PATH
	filePath := filepath.Join(dir, fileName)
	log.Println(filePath)

	// Delete the file if no lines are left
	err = os.Remove(filePath)
	if err != nil {
		http.Error(w, "Unable to delete file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+view, http.StatusSeeOther)
}
