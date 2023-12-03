package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
	HasDirectories bool
	Directories    []GalleryDirectoryData
	Files          []GalleryFileData
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

type RequestHandlers struct {
	MediaDirectory string
	Stat           func(name string) (fs.FileInfo, error)
	ReadDir        func(name string) ([]fs.DirEntry, error)
	Templates      *template.Template
}

func (hdlr RequestHandlers) serveFile(w http.ResponseWriter, r *http.Request, f *os.File) {
	fileInfo, err := f.Stat()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), f)
}

func (hdlr RequestHandlers) getMediaFile(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Path
	filepath = strings.Replace(filepath, "/_media", hdlr.MediaDirectory, 1)
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	defer file.Close()
	hdlr.serveFile(w, r, file)
}

func (hdlr RequestHandlers) getThumbnail(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Path
	filepath = strings.Replace(filepath, "/_thumbnail", hdlr.MediaDirectory, 1)
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	defer file.Close()
	hdlr.serveFile(w, r, file)
}

func (hdlr RequestHandlers) getStaticFile(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Path

	filepath = strings.Replace(filepath, "/static", "./static", 1)
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	defer file.Close()
	hdlr.serveFile(w, r, file)
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

func (hdlr RequestHandlers) getPageData(path string) *PageData {
	requestDir := hdlr.MediaDirectory
	if path != "/" {
		requestDir = requestDir + path
	}

	compoundLink := ""

	breadcrumbs := []Breadcrumb{}

	for _, prt := range strings.Split(path, "/") {
		compoundLink = compoundLink + prt + "/"
		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Name: prt,
			Link: compoundLink,
		})
	}

	data := PageData{
		ShowBreadcrumb: path != "/",
		Breadcrumbs:    breadcrumbs,
		ShowGallery:    true,
	}

	checkFile, err := hdlr.Stat(requestDir)

	if errors.Is(err, os.ErrNotExist) || checkFile.IsDir() {
		files, err := hdlr.ReadDir(requestDir)
		if err != nil {
			return nil
		}
		directories := []GalleryDirectoryData{}
		galleryFiles := []GalleryFileData{}

		rooting := path
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
					Thumbnail: fmt.Sprintf("/_thumbnail%s/%s", path, file.Name()),
				})
			}
		}
		data.GalleryData = &GalleryData{
			HasDirectories: len(directories) > 0,
			Directories:    directories,
			Files:          galleryFiles,
		}
		data.ShowGallery = true
	} else {
		data.ShowGallery = false
		data.ImageData = &ImageData{
			RawPath: "/_media" + path,
		}
	}
	return &data
}

func (hdlr RequestHandlers) handlePage(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		if strings.HasPrefix(request.URL.Path, "/_media") {
			hdlr.getMediaFile(writer, request)
			return
		}
		if strings.HasPrefix(request.URL.Path, "/_thumbnail") {
			hdlr.getThumbnail(writer, request)
			return
		}
		if strings.HasPrefix(request.URL.Path, "/static") {
			hdlr.getStaticFile(writer, request)
			return
		}

		data := hdlr.getPageData(request.URL.Path)

		if data == nil {
			return
		}

		err := hdlr.Templates.ExecuteTemplate(writer, "baseHTML", data)
		if err != nil {
			return
		}
	}
}

func main() {
	gotTemplates, err := getTemplates()
	if err != nil {
		fmt.Printf("error initialising server: %s\n", err)
		os.Exit(1)
	}
	mediaDir := os.Getenv("SMG_MEDIA_DIRECTORY")
	if mediaDir == "" {
		mediaDir = "/_media"
	}
	mux := http.NewServeMux()

	hdlr := RequestHandlers{
		MediaDirectory: mediaDir,
		Stat:           os.Stat,
		ReadDir:        os.ReadDir,
		Templates:      gotTemplates,
	}

	mux.HandleFunc("*", hdlr.handlePage)
	mux.HandleFunc("/", hdlr.handlePage)

	port := "3333"
	portSetting := os.Getenv("SMG_PORT")
	if portSetting != "" {
		port = portSetting
	}
	err = http.ListenAndServe(":"+port, mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
