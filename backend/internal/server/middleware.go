package server

import (
	"net/http"
	"strings"

	"movie-quizzer/backend/internal/service"
)

func (s *Server) auth(next func(w http.ResponseWriter, r *http.Request, user service.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" {
			Error(w, "unauthorized", "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 {
			Error(w, "unauthorized", "invalid Authorization", http.StatusUnauthorized)
			return
		}

		user, err := s.svc.UserFromToken(parts[1])
		if err != nil {
			Error(w, "unauthorized", err.Error(), http.StatusUnauthorized)
			return
		}

		next(w, r, user)
	}
}
