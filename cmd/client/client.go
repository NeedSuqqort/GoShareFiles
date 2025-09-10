package client

import (
	fileHandler "filesharing/pkg/FileServerHandler"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func Init () {
	http.HandleFunc("/uploads/", uploadsHandler)
    http.HandleFunc("/uploadFile/", uploadFileHandler)
}

func render (writer http.ResponseWriter, data interface{}, templates ...string ) {
	t, err := template.ParseFiles(templates...)
	err = t.Execute(writer, data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func uploadsHandler (writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	fmt.Println(path)
	data := fileHandler.GetServerFilesHandler(strings.TrimPrefix(path, fileHandler.UploadPath[1:]))

	if data == nil {
		http.Error(writer, "Could not read the file server. Please try again later", http.StatusInternalServerError)
	}
	render(writer,  data, "./templates/layout.html", "./templates/uploads.html")
}

func uploadFileHandler (writer http.ResponseWriter, request *http.Request) {
	render(writer, nil, "./templates/layout.html", "./templates/uploadFile.html")
}