package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lckrugel/file-server/internal/auth"
	"github.com/lckrugel/file-server/internal/users"
)

func (s *APIServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		headerParts := strings.Split(authHeader, "Bearer ")
		jwt := headerParts[1]

		if authHeader == "" || len(headerParts) != 2 || jwt == "" {
			apiErr := unauthorized("Missing JWT")
			apiErr.write(w)
			return
		}

		claims, err := auth.ValidateJWT(jwt)
		if errors.Is(err, auth.ErrInvalidJWT) {
			apiErr := unauthorized("Invalid JWT")
			apiErr.write(w)
			return
		} else if err != nil {
			log.Printf("failed to validate JWT: %v", err)
			apiErr := internalError("Failed to validate JWT", err)
			apiErr.write(w)
			return
		}

		if claims.ExpiresAt.Before(time.Now()) {
			apiErr := unauthorized("Invalid JWT")
			apiErr.write(w)
			return
		}

		user, err := s.userService.FindById(r.Context(), claims.UserID)
		if errors.Is(err, users.ErrUserNotFound) {
			apiErr := unauthorized("Invalid JWT")
			apiErr.write(w)
			return
		} else if err != nil {
			log.Printf("failed to fetch user: %v", err)
			apiErr := internalError("Failed to fatch user JWT", err)
			apiErr.write(w)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
