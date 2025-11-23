package server

import (
	"movie-quizzer/backend/internal/service"
	"net/http"
)

type Server struct {
	svc service.Service
}

func New(svc service.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// health
	mux.HandleFunc("GET /health", s.health)

	// auth
	mux.HandleFunc("POST /api/auth/register", s.register)
	mux.HandleFunc("POST /api/auth/login", s.login)
	mux.HandleFunc("GET /api/auth/me", s.auth(s.me))

	// game
	mux.HandleFunc("POST /api/game/sessions", s.auth(s.startSession))
	mux.HandleFunc("POST /api/game/sessions/{id}/answers", s.auth(s.submitAnswer))
	mux.HandleFunc("GET /api/game/sessions/{id}/next", s.auth(s.nextQuestion))
	mux.HandleFunc("GET /api/game/sessions/{id}", s.auth(s.sessionSummary))

	// profile
	mux.HandleFunc("GET /api/profile", s.auth(s.profile))
	mux.HandleFunc("GET /api/profile/achievements", s.auth(s.achievements))

	// leaderboard
	mux.HandleFunc("GET /api/leaderboard/global", s.globalLeaderboard)
	mux.HandleFunc("GET /api/leaderboard/me", s.auth(s.myLeaderboard))

	// media
	mux.HandleFunc("GET /api/media/{bucket}/{file}", s.media)

	return mux
}
