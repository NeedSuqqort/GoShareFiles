package fileserverhandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const UploadPath = "./uploads/"
type FileInfo struct {
	Name string `json:"name"`
 	Size int64 	`json:"size"`
	Path string `json:"path"`
	Time string `json:"time"`
	Mode os.FileMode `json:"mode"`
	IsDir bool `json:"isDir"`
}


func Init () {
	os.MkdirAll(UploadPath, os.ModePerm)
	http.HandleFunc("/upload/", uploadHandler)
    http.HandleFunc("/download/", downloadHandler)
}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	
	switch request.Method {

		case http.MethodPost:
			if request.Method != http.MethodPost {
				http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			file, header, err := request.FormFile("file")
			if err != nil {
				http.Error(writer, "Could not receive your uploaded file, please try again.", http.StatusBadRequest)
			}
			defer file.Close()

			outFile, err := os.Create(UploadPath + header.Filename)
			if err != nil {
				http.Error(writer, "Unable to create the file.", http.StatusInternalServerError)
				return
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, file) 
			if err != nil {
				http.Error(writer, "Unable to save the file. ", http.StatusInternalServerError)
			}

			fmt.Fprintf(writer, "File %s uploaded successfully.", header.Filename)

		default:
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		
	}
}

func downloadHandler(writer http.ResponseWriter, request *http.Request) {
	filename := request.URL.Path[len("/download/"):]
	fmt.Println(filename)

	file, err := os.Open(UploadPath + filename)
	if err != nil {
		http.Error(writer, "Unable to locate the file. Check filename.", http.StatusNotFound)
		return
	}
	defer file.Close()
	writer.Header().Set("Content-Disposition", "attachment; filename="+ filename)
    writer.Header().Set("Content-Type", "application/octet-stream")
    io.Copy(writer, file)
	
	fmt.Fprintf(writer, "File %s downloaded successfully.", filename)
}

func GetServerFilesHandler(path string) map[string]interface{} {
	files, err := os.ReadDir(filepath.Join(UploadPath,path))

	if err != nil {
		fmt.Println(err)
		return nil
	}
	var fileInfos []FileInfo

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		fileInfos = append(fileInfos, FileInfo{
			Name: fileInfo.Name(),
			Size: fileInfo.Size(),
			Time: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
			Mode: fileInfo.Mode(),
			Path: filepath.Join(path,fileInfo.Name()),
			IsDir: file.IsDir(),
		})
	}

	return map[string]interface{}{"Files": fileInfos}
}