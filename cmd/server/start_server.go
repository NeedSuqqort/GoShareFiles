package cmd

import (
	fileServer "filesharing/pkg/FileServerHandler"
	"filesharing/pkg/sysinfo"
	"fmt"
	"log"
	"net/http"
)

func Start_server() {
	fmt.Println("Go Version:", sysinfo.GoVersion())

	http.Handle("/", http.FileServer(http.Dir("./static")))
	fileServer.Init()
	// http.HandleFunc("/increment", incrementCounter)

	fmt.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}