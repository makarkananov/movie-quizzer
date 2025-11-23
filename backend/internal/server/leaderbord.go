package server

import (
	"movie-quizzer/backend/internal/service"
	"net/http"
	"strconv"
)

func (s *Server) globalLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}

	list, err := s.svc.GetGlobalLeaderboard(limit)
	if err != nil {
		Error(w, "leaderboard_error", err.Error(), http.StatusInternalServerError)
		return
	}

	JSON(w, list, http.StatusOK)
}

func (s *Server) myLeaderboard(w http.ResponseWriter, r *http.Request, user service.User) {
	entry, err := s.svc.GetLeaderboardEntry(user.ID)
	if err != nil {
		Error(w, "leaderboard_error", err.Error(), http.StatusInternalServerError)
		return
	}

	JSON(w, entry, http.StatusOK)
}
