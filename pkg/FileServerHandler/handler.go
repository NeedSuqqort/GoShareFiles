package fileserverhandler

import (
	"archive/zip"
	"encoding/json"
	db "filesharing/internal/data"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const UploadPath = "./uploads/"

type FileInfo struct {
	Name  string      `json:"name"`
	Size  int64       `json:"size"`
	Path  string      `json:"path"`
	Time  string      `json:"time"`
	Mode  os.FileMode `json:"mode"`
	IsDir bool        `json:"isDir"`
}

type RepoName struct {
    Name string `json:"name"`
}

type Response struct {
    AccessCode string `json:"access_code"`
    Status     int    `json:"status"`
}

func Init() {
	os.MkdirAll(UploadPath, os.ModePerm)
	http.HandleFunc("/upload/", uploadHandler)
	http.HandleFunc("/files/", filesHandler)
	http.HandleFunc("/create-folder/", folderHandler)
	http.HandleFunc("/create-repo/", repoHandler)
}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {


	switch request.Method {

		case "POST":

			file, header, err := request.FormFile("file")
			path := request.FormValue("path")
			fmt.Println("Upload path:", path)

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

func filesHandler(writer http.ResponseWriter, request *http.Request) {
	filename := strings.TrimPrefix(request.URL.Path, "/files/")

	switch request.Method {

		case "GET":
			_, e := os.ReadDir(filepath.Join(UploadPath, filename))
			if e != nil {
				if pathErr, ok := e.(*os.PathError); ok {
					file, err := os.Open(filepath.Join(UploadPath, filename))
					if err != nil {
						http.Error(writer, "Unable to read the file. Check if the file is corrupted.", http.StatusNotFound)
						return
					}
					defer file.Close()
					writer.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filename))
					writer.Header().Set("Content-Type", "application/octet-stream")
					io.Copy(writer, file)
					fmt.Fprintf(writer, "File %s downloaded successfully.", filename)

				} else {
					http.Error(writer, pathErr.Error(), http.StatusNotFound)
				}
			} else {
				
				zipfile := filepath.Base(filename) + ".zip"
				zipWriter := zip.NewWriter(writer)

				defer zipWriter.Close()
								
				writer.Header().Set("Content-Disposition", "attachment; filename="+zipfile)
				writer.Header().Set("Content-Type", "application/zip")

				err := filepath.Walk(filepath.Join(UploadPath,filename), func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if info.IsDir() {
						return nil
					}

					path = filepath.ToSlash(path)
					file, err := os.Open(path)
					if err != nil {
						return err
					}
					defer file.Close()

					fmt.Println("Zipping file:", path)
					fmt.Println("Copying from:", file.Name())

					zipFile, err := zipWriter.Create(path)
					if err != nil {
						return err
					}

					_, err = io.Copy(zipFile, file)
					return err
				})

				if err != nil {
					http.Error(writer, "Unable to zip the folder. Please retry by downloading again.", http.StatusInternalServerError)
				}
			}
		
		case "POST":

			err := request.ParseMultipartForm(200 << 20) // Max 200MB
			
			if err != nil {
				http.Error(writer, "Wrong format.", http.StatusBadRequest)
			}

			path := request.FormValue("path")
			fmt.Println("Upload path:", path)
			file, header, err := request.FormFile("file")

			if err != nil {
				http.Error(writer, "Could not receive your uploaded file, please try again.", http.StatusBadRequest)
			}
			defer file.Close()

			outFile, err := os.Create(UploadPath + strings.TrimPrefix(path, UploadPath[1:]) + "/" + header.Filename)
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

		case "DELETE":
			target := UploadPath + filename
			err := os.RemoveAll(target)

			if err != nil {
				if os.IsNotExist(err) {
					info, statErr := os.Stat(target) 
					if statErr == nil && info.IsDir() {
						http.Error(writer, "Folder does not exist.", http.StatusNotFound)
					} else {
						http.Error(writer, "File does not exist.", http.StatusNotFound)
					}
				} else {
					http.Error(writer, "Failed to delete the directory, please try again later.", http.StatusInternalServerError)
				}
			}

			fmt.Fprintf(writer, "File %s downloaded successfully.", filename)

		default:
			http.Error(writer, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}

func GetServerFilesHandler(path string) map[string]interface{} {
	files, err := os.ReadDir(filepath.Join(UploadPath, path))

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
			Name:  fileInfo.Name(),
			Size:  fileInfo.Size(),
			Time:  fileInfo.ModTime().Format("2006-01-02 15:04:05"),
			Mode:  fileInfo.Mode(),
			Path:  filepath.Join(path, fileInfo.Name()),
			IsDir: file.IsDir(),
		})
	}

	return map[string]interface{}{"Files": fileInfos}
}

func folderHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.URL.Path)
	targetPath := strings.TrimPrefix(request.URL.Path, "/create-folder/")
	_, err := os.Stat(targetPath)

	if os.IsNotExist(err) {
		err := os.Mkdir(targetPath, 0755)
		if err != nil {
			fmt.Println(err)
			http.Error(writer, "Failed to create the folder, please refresh and try again.", http.StatusInternalServerError)
		}
		fmt.Println("Success")
	} else {
		http.Error(writer, "Folder already exists, please use a different name.", http.StatusInternalServerError)
	}
}

func repoHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var name RepoName

	if e := json.NewDecoder(request.Body).Decode(&name); e != nil {
        http.Error(writer, e.Error(), http.StatusBadRequest)
        return
    }

	access_code, e := db.AddRepo(name.Name)
	
	if e != nil {
		http.Error(writer, e.Error(), http.StatusInternalServerError)
		return
	}

	targetPath := UploadPath + access_code

	_, err := os.Stat(targetPath)

	if os.IsNotExist(err) {
		err := os.Mkdir(targetPath, 0755)

		if err != nil {
			fmt.Println(err)
			http.Error(writer, "Failed to create the repo, please refresh and try again.", http.StatusInternalServerError)
		}
		fmt.Println("Success")

		writer.Header().Set("Content-Type", "application/json")
		response := Response{
			AccessCode: access_code,
			Status: 200,
		}

		if err := json.NewEncoder(writer).Encode(response); err != nil {
        	http.Error(writer, err.Error(), http.StatusInternalServerError)
    	}

	} else {
		http.Error(writer, "Repo already exists, please use a different name.", http.StatusInternalServerError)
	}
}
