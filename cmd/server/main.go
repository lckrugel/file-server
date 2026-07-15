package main

import (
	"log"
	"net/http"

	"github.com/lckrugel/file-server/internal/auth"
	"github.com/lckrugel/file-server/internal/files"
	"github.com/lckrugel/file-server/internal/server"
	"github.com/lckrugel/file-server/internal/users"
)

func main() {
	fileRepo := files.NewMemoryRepository()
	fileStorage := files.NewLocalStorage("./storage")
	fileService := files.NewFileService(fileRepo, fileStorage)

	userRepo := users.NewMemoryRepository()
	userService := users.NewUserService(userRepo)

	credentialRepo := auth.NewMemoryRepository()
	authService := auth.NewAuthService(credentialRepo, userService)

	server := server.NewAPIServer(
		authService,
		userService,
		fileService,
	)

	router := server.Router()

	log.Println("Starting server on localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
