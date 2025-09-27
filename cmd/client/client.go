package client

import (
	fileHandler "filesharing/pkg/FileServerHandler"
	"filesharing/pkg/middleware"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func Init() {
	http.Handle("/uploads/",middleware.BreadcrumbMiddleware(http.HandlerFunc(uploadsHandler)))
    // http.HandleFunc("/uploadFile/", uploadFileHandler)
	http.HandleFunc("/", mainPageHandler)
}

func render(writer http.ResponseWriter, data interface{}, templates ...string) {
	t := template.Must(template.ParseFiles(templates...))
	err := t.Execute(writer, data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadsHandler(writer http.ResponseWriter, request *http.Request) {
	
	path := strings.TrimPrefix(request.URL.Path, fileHandler.UploadPath[1:])
	if path == "" {
		http.Redirect(writer, request, "/", http.StatusMovedPermanently)
		return
	}

	data := fileHandler.GetServerFilesHandler(path)
	if data == nil {
		http.Error(writer, "Cannot locate the resource.", http.StatusNotFound)
		return
	}

	breadcrumbs := middleware.GetBreadcrumbs(request)
	data["Breadcrumbs"] = breadcrumbs
	fmt.Println(path)
	data["AccessCode"] = path[:5]
	render(writer,  data, "./templates/layout.html", "./templates/breadcrumb.html", "./templates/uploads.html")
}

func uploadFileHandler(writer http.ResponseWriter, request *http.Request) {
	render(writer, nil, "./templates/layout.html", "./templates/uploadFile.html")
}

func mainPageHandler(writer http.ResponseWriter, request *http.Request) {
	render(writer, nil, "./templates/layout.html", "./templates/index.html")
}	
