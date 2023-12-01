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

var mediaDir string

func serveFile(w http.ResponseWriter, r *http.Request, f *os.File) {
	fileInfo, err := f.Stat()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), f)
}

func getMediaFile(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Path
	filepath = strings.Replace(filepath, "/media", mediaDir, 1)
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	defer file.Close()
	serveFile(w, r, file)
}

func getThumbnail(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Path
	mediaDir := os.Getenv("SMG_MEDIA_DIRECTORY")
	if mediaDir == "" {
		mediaDir = "/_media"
	}
	filepath = strings.Replace(filepath, "/thumbnail", mediaDir, 1)
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	defer file.Close()
	serveFile(w, r, file)
}

func getStaticFile(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Path

	filepath = strings.Replace(filepath, "/static", "./static", 1)
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	defer file.Close()
	serveFile(w, r, file)
}

func getTemplates() (templates *template.Template, err error) {
	var allFiles []string
	files, _ := os.ReadDir("templates")
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".gohtml") {
			filePath := filepath.Join("templates", filename)
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
	mediaDir = os.Getenv("SMG_MEDIA_DIRECTORY")
	if mediaDir == "" {
		mediaDir = "/_media"
	}
}

type GalleryDirectoryData struct {
	Name string
	Link string
}

type GalleryFileData struct {
	Name      string
	Link      string
	Thumbnail string
}

type GalleryData struct {
	Directories []GalleryDirectoryData
	Files       []GalleryFileData
}

type ImageData struct {
	RawPath string
}

type Breadcrumb struct {
	Name string
	Link string
}

type PageData struct {
	ShowBreadcrumb bool
	Breadcrumbs    []Breadcrumb
	ShowGallery    bool
	GalleryData    *GalleryData
	ImageData      *ImageData
}

func handlePage(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		if strings.HasPrefix(request.URL.Path, "/media") {
			getMediaFile(writer, request)
			return
		}
		if strings.HasPrefix(request.URL.Path, "/thumbnail") {
			getThumbnail(writer, request)
			return
		}
		if strings.HasPrefix(request.URL.Path, "/static") {
			getStaticFile(writer, request)
			return
		}

		requestDir := mediaDir
		if request.URL.Path != "/" {
			requestDir = requestDir + request.URL.Path
		}

		compoundLink := ""

		breadcrumbs := []Breadcrumb{}

		for _, prt := range strings.Split(request.URL.Path, "/") {
			compoundLink = compoundLink + prt + "/"
			breadcrumbs = append(breadcrumbs, Breadcrumb{
				Name: prt,
				Link: compoundLink,
			})
		}

		data := PageData{
			ShowBreadcrumb: request.URL.Path != "/",
			Breadcrumbs:    breadcrumbs,
			ShowGallery:    true,
		}

		checkFile, err := os.Stat(requestDir)

		if errors.Is(err, os.ErrNotExist) || checkFile.IsDir() {
			files, err := os.ReadDir(requestDir)
			if err != nil {
				return // 404
			}
			directories := []GalleryDirectoryData{}
			galleryFiles := []GalleryFileData{}

			rooting := request.URL.Path
			if rooting == "/" {
				rooting = ""
			}
			for _, file := range files {
				if file.IsDir() {
					directories = append(directories, GalleryDirectoryData{
						Name: file.Name(),
						Link: fmt.Sprintf("%s/%s", rooting, file.Name()),
					})
				} else {
					galleryFiles = append(galleryFiles, GalleryFileData{
						Name:      file.Name(),
						Link:      fmt.Sprintf("%s/%s", rooting, file.Name()),
						Thumbnail: fmt.Sprintf("/thumbnail%s/%s", request.URL.Path, file.Name()),
					})
				}
			}
			data.GalleryData = &GalleryData{
				Directories: directories,
				Files:       galleryFiles,
			}
			data.ShowGallery = true
		} else {
			data.ShowGallery = false
			data.ImageData = &ImageData{
				RawPath: "/media" + request.URL.Path,
			}
		}

		err = templates.ExecuteTemplate(writer, "baseHTML", data)

		if err != nil {
			return
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("*", handlePage)
	mux.HandleFunc("/", handlePage)

	port := "3333"
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
