package server

import (
	"movie-quizzer/backend/internal/service"
	"net/http"
)

func (s *Server) profile(w http.ResponseWriter, r *http.Request, user service.User) {
	p, err := s.svc.GetProfile(user.ID)
	if err != nil {
		Error(w, "profile_error", err.Error(), http.StatusInternalServerError)
		return
	}

	JSON(w, p, http.StatusOK)
}

func (s *Server) achievements(w http.ResponseWriter, r *http.Request, user service.User) {
	a, err := s.svc.GetAchievements(user.ID)
	if err != nil {
		Error(w, "achievements_error", err.Error(), http.StatusInternalServerError)
		return
	}

	JSON(w, a, http.StatusOK)
}
