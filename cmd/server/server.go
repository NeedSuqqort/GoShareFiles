package cmd

import (
	client "filesharing/cmd/client"
	db "filesharing/internal/data"
	server "filesharing/pkg/FileServerHandler"
	"filesharing/pkg/sysinfo"
	"fmt"
	"log"
	"net/http"
)

func Start_server() {
	fmt.Println("Go Version:", sysinfo.GoVersion())
	client.Init()
	server.Init()
	db.Init()
	fmt.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}