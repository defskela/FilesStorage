package main

import (
	"filesStorage/upload"

	"github.com/defskela/httpServer/router"
	"github.com/defskela/httpServer/server"
)

func main() {
	router := router.NewRouter()
	router.Post("/upload", upload.UploadHandler)
	server.StartServ(router, 0)
}
