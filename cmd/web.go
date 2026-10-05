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
	"time"

	"github.com/radoslawg/video_manager/resources"
	"github.com/spf13/cobra"
)

var port int16 = 8080
var address string = ""
var videosPath = "z:/youtube/!nextdaily"

var templates *template.Template

func init() {
	webCmd.Flags().Int16VarP(&port, "port", "p", 8080, "Port number for web server")
	webCmd.Flags().StringVarP(&address, "bind", "b", address, "Bind to address")
	webCmd.Flags().StringVarP(&videosPath, "videos_path", "v", videosPath, "Path to videos")
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
		server.HandleFunc("/delete-all/", deleteAllHandler)
		server.Handle("/static/", http.FileServer(http.FS(resources.StaticFiles)))

		fmt.Printf("Starting Web server on %v:%v\n", address, port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), server)
		if err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	},
}

type FileLinks struct {
	FileName         string
	Links            []string
	OriginalFileName []string
	Titles           []string
	Dates            []string
	Times            []string
}

//const videosPath = "/media/sda1/youtube/!nextdaily"

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	dir := videosPath // Directory containing the files
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mhtml") {
			dates := strings.Split(strings.TrimSpace(string(file.Name())), "#")
			day := dateOnly(dates[0])
			if !slices.Contains(fileNames, day) {
				fileNames = append(fileNames, day)
			}
		}
	}

	templates.Lookup("index.tmpl").Execute(w, fileNames)
}

// dateOnly truncates a 14-digit dayKey (YYYYmmddhhMMss) to just its
// 8-digit date portion (YYYYmmdd), so videos downloaded on the same
// day at different times group together. 8-digit dayKeys pass through
// unchanged.
func dateOnly(dayKey string) string {
	if len(dayKey) == 14 {
		return dayKey[:8]
	}
	return dayKey
}

// formatDate renders the grouping date of a dayKey (YYYYmmdd or
// YYYYmmddhhMMss) as "YYYY-MM-DD", or "-" if it can't be parsed as a date.
func formatDate(dayKey string) string {
	d := dateOnly(dayKey)
	t, err := time.Parse("20060102", d)
	if err != nil {
		return "-"
	}
	return t.Format("2006-01-02")
}

// formatTime renders the time-of-day portion of a 14-digit dayKey
// (YYYYmmddhhMMss) as "HH:MM:SS", or "-" if the dayKey has no time
// portion (8-digit dayKeys) or can't be parsed.
func formatTime(dayKey string) string {
	if len(dayKey) != 14 {
		return "-"
	}
	t, err := time.Parse("20060102150405", dayKey)
	if err != nil {
		return "-"
	}
	return t.Format("15:04:05")
}

func viewFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/view/")
	fileName = filepath.Base(fileName) // Prevent directory traversal
	dir := videosPath                  // Directory containing the files

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	var links []string
	var original_filenames []string
	var titles []string
	var dates []string
	var times []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mhtml") && strings.HasPrefix(file.Name(), fileName) {
			i := 2
			ids := strings.Split(strings.TrimSpace(string(file.Name())), "#")
			if len(ids[1]) == 8 || len(ids[1]) == 14 {
				_, err := strconv.ParseInt(ids[1], 10, 64)
				if err == nil {
					i = 3
				}
			}
			if !slices.Contains(links, ids[i]) {
				original_filenames = append(original_filenames, file.Name())
				links = append(links, ids[i])
				titles = append(titles, strings.ReplaceAll(ids[i-1], "_", " "))
				dates = append(dates, formatDate(ids[0]))
				times = append(times, formatTime(ids[0]))
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
		Dates:            dates,
		Times:            times,
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

	fileName, err := url.PathUnescape(parts[0]) // Prevent directory traversal
	if err != nil {
		http.Error(w, "Unable to parse filename", http.StatusInternalServerError)
		return
	}
	view := parts[1]
	dir := videosPath
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

func deleteAllHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/delete-all/")
	fileName = filepath.Base(fileName) // Prevent directory traversal
	dir := videosPath

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".mhtml") && strings.HasPrefix(file.Name(), fileName) {
			filePath := filepath.Join(dir, file.Name())
			if err := os.Remove(filePath); err != nil {
				log.Println(err)
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
