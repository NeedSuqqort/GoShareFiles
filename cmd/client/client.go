package client

import (
	"html/template"
	"net/http"
)


func Init () {
	http.HandleFunc("/uploads/", uploadsHandler)
    http.HandleFunc("/uploadFile/", uploadFileHandler)
}

func render (writer http.ResponseWriter, templates ...string) {
	t, err := template.ParseFiles(templates...)
	err = t.Execute(writer, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func uploadsHandler (writer http.ResponseWriter, request *http.Request) {
	render(writer, "./templates/layout.html", "./templates/uploads.html")
}

func uploadFileHandler (writer http.ResponseWriter, request *http.Request) {
	render(writer, "./templates/layout.html", "./templates/uploadFile.html")
}