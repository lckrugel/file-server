package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type authResponse struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func (s *APIServer) register(w http.ResponseWriter, req *http.Request) {
	type RegistrationBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	var data RegistrationBody

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		apiErr := badRequest("invalid request body")
		apiErr.write(w)
	}

	if parts := strings.Split(data.Email, "@"); len(parts) < 2 {
		apiErr := badRequest("invalid email")
		apiErr.write(w)
	}

	user, jwt, err := s.authService.Register(req.Context(), "", "", "")
	if err != nil {
		log.Printf("failed to register user: %v", err)
		apiErr := internalError("failed to get file", err)
		apiErr.write(w)
		return
	}

	type RegistrationResponse struct {
		UserID  string        `json:"user_id"`
		Email   string        `json:"email"`
		Name    string        `json:"name"`
		Session *authResponse `json:"session"`
	}
	resp := created("New user registered", &RegistrationResponse{
		UserID: user.ID.String(),
		Email:  user.Email,
		Name:   user.Name,
		Session: &authResponse{
			Type:  "Bearer",
			Token: jwt,
		},
	})
	resp.write(w)
}

// func (s *APIServer) Login(w http.ResponseWriter, req *http.Request) {}
