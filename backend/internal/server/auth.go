package server

import (
	"encoding/json"
	"net/http"

	"movie-quizzer/backend/internal/service"
)

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRes struct {
	Token string `json:"token"`
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	JSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, "invalid_json", "bad request", http.StatusBadRequest)
		return
	}

	if err := s.svc.Register(req.Email, req.Password, req.Nickname); err != nil {
		Error(w, "register_failed", err.Error(), http.StatusBadRequest)
		return
	}

	JSON(w, map[string]string{"status": "registered"}, http.StatusCreated)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, "invalid_json", "bad request", http.StatusBadRequest)
		return
	}

	token, err := s.svc.Login(req.Email, req.Password)
	if err != nil {
		Error(w, "login_failed", err.Error(), http.StatusUnauthorized)
		return
	}

	JSON(w, loginRes{Token: token}, http.StatusOK)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request, user service.User) {
	JSON(w, user, http.StatusOK)
}
