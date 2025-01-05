package main

import (
	"filesStorage/upload"

	"github.com/defskela/httpServer/router"
	httpServer "github.com/defskela/httpServer/server"
)

func main() {
	router := router.NewRouter()
	router.Post("/upload", upload.UploadHandler)
	server := httpServer.NewServer(router)
	server.Start("8080")
}
