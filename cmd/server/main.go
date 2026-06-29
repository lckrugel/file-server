package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lckrugel/file-server/internal/files"
	"github.com/lckrugel/file-server/internal/handlers"
)

func main() {
	router := chi.NewRouter()

	fileRepo := files.NewMemoryRepository()
	fileStorage := files.NewLocalStorage("./storage")
	fileService := files.NewFileService(fileRepo, fileStorage)
	fileHandler := handlers.NewFileHandler(fileService)

	router.Post("/files", fileHandler.Upload)
	router.Get("/files/{fileId}", fileHandler.Download)
	router.Get("/files", fileHandler.List)

	log.Println("Starting server on localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
