package fileserverhandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const uploadPath = "./uploads/"

func Init () {
	os.MkdirAll(uploadPath, os.ModePerm)
	http.HandleFunc("/upload", uploadHandler)
    http.HandleFunc("/download/", downloadHandler)
}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := request.FormFile("file")
	if err != nil {
		http.Error(writer, "Could not receive your uploaded file, please try again.", http.StatusBadRequest)
	}
	defer file.Close()

	outFile, err := os.Create(uploadPath + header.Filename)
	if err != nil {
		http.Error(writer, "Unable to create the file. ", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file) 
	if err != nil {
		http.Error(writer, "Unable to save the file. ", http.StatusInternalServerError)
	}

	fmt.Fprintf(writer, "File %s uploaded successfully.", header.Filename)

}

func downloadHandler(writer http.ResponseWriter, request *http.Request) {
	filename := request.URL.Path[len("/download/"):]

	file, err := os.Open(uploadPath + filename)
	if err != nil {
		http.Error(writer, "Unable to locate the file. Check filename.", http.StatusNotFound)
	}
	defer file.Close()
	writer.Header().Set("Content-Disposition", "attachment; filename="+ filename)
    writer.Header().Set("Content-Type", "application/octet-stream")
    io.Copy(writer, file)
	
	fmt.Fprintf(writer, "File %s downloaded successfully.", filename)
}