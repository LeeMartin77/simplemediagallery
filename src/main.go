package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template

func getTemplates() (templates *template.Template, err error) {
	var allFiles []string
	files2, _ := os.ReadDir("_templates")
	for _, file := range files2 {
		filename := file.Name()
		if strings.HasSuffix(filename, ".gohtml") {
			filePath := filepath.Join("_templates", filename)
			allFiles = append(allFiles, filePath)
		}
	}

	return template.New("").ParseFiles(allFiles...)
}

func init() {
	gotTemplates, err := getTemplates()
	if err != nil {
		fmt.Printf("error initialising server: %s\n", err)
		os.Exit(1)
	}
	templates = gotTemplates
}

func handlePage(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		// event := News{
		// 		Headline: "makeuseof.com has everything Tech",
		// 		Body: "Visit MUO for anything technology related",
		// }

		data := 2

		err := templates.ExecuteTemplate(writer, "baseHTML", data)

		if err != nil {
			return
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("*", handlePage)
	mux.HandleFunc("/", handlePage)

	port := "8080"
	portSetting := os.Getenv("SMG_PORT")
	if portSetting != "" {
		port = portSetting
	}
	err := http.ListenAndServe(":"+port, mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
