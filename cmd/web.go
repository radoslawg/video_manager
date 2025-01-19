package cmd

import (
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

	"github.com/radoslawg/video_manager/resources"
	"github.com/spf13/cobra"
)

var port int16 = 8080
var address string = ""

var templates *template.Template

func init() {
	webCmd.Flags().Int16VarP(&port, "port", "p", 8080, "Port number for web server")
	webCmd.Flags().StringVarP(&address, "bind", "b", address, "Bind to address")
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start Web Server",
	Long:  `Starts Web Server to access configuration and player of video-manager`,
	Run: func(cmd *cobra.Command, args []string) {
		templates = resources.Templates()

		server := http.NewServeMux()
		server.HandleFunc("/", listFilesHandler)
		server.HandleFunc("/view/", viewFileHandler)
		server.HandleFunc("/delete/", deleteLinkHandler)
		server.Handle("/static/", http.FileServer(http.FS(resources.StaticFiles)))

		fmt.Printf("Starting Web server on %v:%v\n", address, port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), server)
	},
}

type FileLinks struct {
	FileName         string
	Links            []string
	OriginalFileName []string
	Titles           []string
}

//const PATH = "/media/sda1/youtube/!nextdaily"

const PATH = "z:/youtube/!nextdaily"

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

	templates.Lookup("index.tmpl").Execute(w, fileNames)
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
	tmpl := templates.Lookup("day_view.tmpl")
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
