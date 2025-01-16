package cmd

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/radoslawg/video_manager/resources"
	"github.com/spf13/cobra"
)

var port int16 = 8080
var address string = ""

func init() {
	webCmd.Flags().Int16VarP(&port, "port", "p", 8080, "Port number for web server")
	webCmd.Flags().StringVarP(&address, "bind", "b", address, "Bind to address")
	rootCmd.AddCommand(webCmd)
}

var templates *template.Template

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start Web Server",
	Long:  `Starts Web Server to access configuration and player of video-manager`,
	Run: func(cmd *cobra.Command, args []string) {
		templates = resources.Templates()

		server := http.NewServeMux()
		server.HandleFunc("/", serveIndex)
		server.Handle("/static/", http.FileServer(http.FS(resources.StaticFiles)))
		fmt.Printf("Starting Web server on %v:%v\n", address, port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), server)
	},
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	err := templates.Execute(w, "index.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Cannot serve template %v", err))
	}
}
