package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/lckrugel/file-server/internal/auth"
	"github.com/lckrugel/file-server/internal/files"
	"github.com/lckrugel/file-server/internal/users"
)

type APIServer struct {
	authService *auth.AuthService
	userService *users.UserService
	fileService *files.FileService
}

func NewAPIServer(
	auth *auth.AuthService,
	user *users.UserService,
	files *files.FileService,
) *APIServer {
	return &APIServer{
		authService: auth,
		userService: user,
		fileService: files,
	}
}

func (s *APIServer) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", s.register)

	r.Group(func(r chi.Router) {
		r.Post("/files", s.upload)
		r.Get("/files/{fileId}", s.download)
		r.Get("/files", s.list)

		r.Use(s.authMiddleware)
	})

	return r
}
