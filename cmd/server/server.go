package cmd

import (
	client "filesharing/cmd/client"
	server "filesharing/pkg/FileServerHandler"
	"filesharing/pkg/sysinfo"
	"fmt"
	"log"
	"net/http"
)

func Start_server() {
	fmt.Println("Go Version:", sysinfo.GoVersion())
	server.Init()
	client.Init()
	fmt.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}