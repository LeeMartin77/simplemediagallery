package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/gabriel-vasile/mimetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
)

type GalleryDirectoryData struct {
	Name      string
	Link      string
	FileCount int
}

type GalleryFileData struct {
	Name, Link, Thumbnail string
}

type GalleryData struct {
	HasDirectories bool
	Directories    []GalleryDirectoryData
	Files          []GalleryFileData
	VisibleTypes   []string
	AvailableTypes []string
}

type FileData struct {
	RawPath  string
	IsImage  bool
	IsVideo  bool
	FileType string
}

type Breadcrumb struct {
	Name, Link string
}

type PageData struct {
	ShowBreadcrumb bool
	Breadcrumbs    []Breadcrumb
	ShowGallery    bool
	GalleryData    *GalleryData
	FileData       *FileData
}

type RequestHandlers struct {
	MediaDirectory string
	Stat           func(name string) (fs.FileInfo, error)
	ReadDir        func(name string) ([]fs.DirEntry, error)
	Templates      *template.Template
	DetectFile     func(path string) (*mimetype.MIME, error)
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

var imageExtensions []string = []string{
	"jpg", "jpeg", "png", "gif", "bmp", "tiff", "webp", "svg",
}

var videoExtensions []string = []string{
	"mp4", "avi", "mkv", "mov", "wmv", "flv", "webm", "mpeg", "mpg", "3gp",
}

func (hdlr RequestHandlers) getThumbnail(w http.ResponseWriter, r *http.Request) {
	filepath := r.URL.Path
	rawWidth := r.URL.Query().Get("width")
	var width uint
	iwidth, err := strconv.Atoi(rawWidth)
	if rawWidth == "" || err != nil {
		iwidth = 300
	}
	width = uint(iwidth)
	filepath = strings.Replace(filepath, "/_thumbnail", hdlr.MediaDirectory, 1)
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	prts := strings.Split(filepath, ".")
	ext := strings.ToLower(prts[len(prts)-1])
	defer file.Close()
	if slices.Contains(imageExtensions, ext) {
		var img image.Image
		var err error
		fileContents, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		mtype := mimetype.Detect(fileContents)
		if !strings.HasPrefix(mtype.String(), "image/") {
			fmt.Println("Not recognised as image")
			return
		}

		imgFormat := strings.Split(mtype.String(), "/")[1]

		switch imgFormat {
		case "jpeg":
			img, err = jpeg.Decode(bytes.NewReader(fileContents))
		case "bmp":
			img, err = bmp.Decode(bytes.NewReader(fileContents))
		case "png":
			img, err = png.Decode(bytes.NewReader(fileContents))
		case "gif":
			img, err = gif.Decode(bytes.NewReader(fileContents))
		// Add more cases for other image formats if needed
		default:
			err = fmt.Errorf("unsupported image format: %s", mtype.String())
		}
		if err != nil {
			return
		}
		st, err := file.Stat()

		imgResized := resize.Resize(width, 0, img, resize.Bilinear)
		var buf bytes.Buffer
		var imgErr error
		switch imgFormat {
		case "jpeg":
			imgErr = jpeg.Encode(&buf, imgResized, nil)
		case "bmp":
			imgErr = bmp.Encode(&buf, imgResized)
		case "png":
			imgErr = png.Encode(&buf, imgResized)
		case "gif":
			imgErr = gif.Encode(&buf, imgResized, nil)
		// Add more cases for other image formats if needed
		default:
			imgErr = fmt.Errorf("unsupported image format: %s", mtype.String())
		}
		if imgErr != nil {
			return
		}
		http.ServeContent(w, r, st.Name(), st.ModTime(), bytes.NewReader(buf.Bytes()))
		return
	}
	if slices.Contains(videoExtensions, ext) {
		file, err := os.Open("./static/play.png")
		if err != nil {
			http.Error(w, "No file found", http.StatusNotFound)
			return
		}
		defer file.Close()
		hdlr.serveFile(w, r, file)
		return
	}
	file, err = os.Open("./static/picture.png")
	if err != nil {
		http.Error(w, "No file found", http.StatusNotFound)
		return
	}
	defer file.Close()
	hdlr.serveFile(w, r, file)
	return
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

func (hdlr RequestHandlers) getPageData(path string, query url.Values) *PageData {
	requestDir := hdlr.MediaDirectory
	breadcrumbs := []Breadcrumb{}
	if path != "/" {
		requestDir = requestDir + path

		compoundLink := ""

		for _, prt := range strings.Split(path, "/") {
			compoundLink = compoundLink + prt + "/"
			breadcrumbs = append(breadcrumbs, Breadcrumb{
				Name: prt,
				Link: compoundLink,
			})
		}

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

		visibleTypes := []string{}
		availableMap := map[string]bool{}

		visible := query["visible"]

		isFiltering := len(visible) > 0

		for _, file := range files {
			if file.IsDir() {
				subdir, err := hdlr.ReadDir(fmt.Sprintf("%s/%s", requestDir, file.Name()))
				if err != nil {
					continue
				}
				directories = append(directories, GalleryDirectoryData{
					Name:      file.Name(),
					Link:      fmt.Sprintf("%s/%s", rooting, file.Name()),
					FileCount: len(subdir),
				})
				continue
			}
			prts := strings.Split(file.Name(), ".")
			ext := prts[len(prts)-1]
			if slices.Contains(imageExtensions, ext) {
				availableMap["image"] = true
				if isFiltering && !slices.Contains(visible, "image") {
					continue
				}
			}
			if slices.Contains(videoExtensions, ext) {
				availableMap["video"] = true
				if isFiltering && !slices.Contains(visible, "video") {
					continue
				}
			}
			if !slices.Contains(imageExtensions, ext) && !slices.Contains(videoExtensions, ext) {
				availableMap["other"] = true
				if isFiltering && !slices.Contains(visible, "other") {
					continue
				}
			}
			extraSlash := ""
			if path != "/" {
				extraSlash = "/"
			}
			galleryFiles = append(galleryFiles, GalleryFileData{
				Name:      file.Name(),
				Link:      fmt.Sprintf("%s/%s", rooting, file.Name()),
				Thumbnail: fmt.Sprintf("/_thumbnail%s%s%s", path, extraSlash, file.Name()),
			})
		}

		availableTypes := []string{}
		for k, v := range availableMap {
			if v {
				availableTypes = append(availableTypes, k)
			}
		}
		sort.SliceStable(availableTypes, func(i, j int) bool {
			return availableTypes[i] > availableTypes[j]
		})
		data.GalleryData = &GalleryData{
			HasDirectories: len(directories) > 0,
			Directories:    directories,
			Files:          galleryFiles,
			VisibleTypes:   visibleTypes,
			AvailableTypes: availableTypes,
		}
		data.ShowGallery = true
	} else {
		data.ShowGallery = false
		mtype, _ := hdlr.DetectFile(requestDir)
		ftype := mtype.String()
		data.FileData = &FileData{
			RawPath:  "/_media" + path,
			IsImage:  strings.HasPrefix(ftype, "image"),
			IsVideo:  strings.HasPrefix(ftype, "video"),
			FileType: ftype,
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

		data := hdlr.getPageData(request.URL.Path, request.URL.Query())

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
	port := "3333"
	portSetting := os.Getenv("SMG_PORT")
	if portSetting != "" {
		port = portSetting
	}
	mux := http.NewServeMux()

	hdlr := RequestHandlers{
		MediaDirectory: mediaDir,
		Stat:           os.Stat,
		ReadDir:        os.ReadDir,
		Templates:      gotTemplates,
		DetectFile:     mimetype.DetectFile,
	}

	mux.HandleFunc("*", hdlr.handlePage)
	mux.HandleFunc("/", hdlr.handlePage)

	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGTERM, syscall.SIGINT)

	srv := &http.Server{Addr: ":" + port, Handler: mux}

	go func() {
		s := <-osSig
		// cleanup
		fmt.Printf("Signal: %v\n", s)
		srv.Shutdown(context.TODO())
	}()

	err = srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
		os.Exit(0)
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
