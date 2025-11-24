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

	styles := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", styles))

	fmt.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}