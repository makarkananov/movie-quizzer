package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"movie-quizzer/backend/internal/service"
)

type startReq struct {
	Mode string `json:"mode"`
}

type answerReq struct {
	QuestionID int64  `json:"question_id"`
	Answer     string `json:"answer"`
	ElapsedMs  int64  `json:"elapsed_ms"`
}

func (s *Server) startSession(w http.ResponseWriter, r *http.Request, user service.User) {
	var req startReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, "invalid_json", "bad request", http.StatusBadRequest)
		return
	}

	sess, q, err := s.svc.StartSession(user.ID, req.Mode)
	if err != nil {
		Error(w, "start_failed", err.Error(), http.StatusBadRequest)
		return
	}

	JSON(w, map[string]any{
		"session":  sess,
		"question": q,
	}, http.StatusCreated)
}

func (s *Server) submitAnswer(w http.ResponseWriter, r *http.Request, user service.User) {
	idStr := r.PathValue("id")
	sID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(w, "invalid_id", "bad session id", http.StatusBadRequest)
		return
	}

	var req answerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, "invalid_json", "bad request", http.StatusBadRequest)
		return
	}

	res, err := s.svc.SubmitAnswer(user.ID, sID, service.SubmitAnswer{
		QuestionID: req.QuestionID,
		Answer:     req.Answer,
		ElapsedMs:  req.ElapsedMs,
	})
	if err != nil {
		Error(w, "submit_failed", err.Error(), http.StatusBadRequest)
		return
	}

	JSON(w, res, http.StatusOK)
}

func (s *Server) nextQuestion(w http.ResponseWriter, r *http.Request, user service.User) {
	idStr := r.PathValue("id")
	sID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(w, "invalid_id", "bad session id", http.StatusBadRequest)
		return
	}

	q, err := s.svc.GetNextQuestion(user.ID, sID)
	if err != nil {
		Error(w, "no_question", err.Error(), http.StatusNotFound)
		return
	}

	JSON(w, q, http.StatusOK)
}

func (s *Server) sessionSummary(w http.ResponseWriter, r *http.Request, user service.User) {
	idStr := r.PathValue("id")
	sID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(w, "invalid_id", "bad session id", http.StatusBadRequest)
		return
	}

	summary, err := s.svc.GetSessionSummary(user.ID, sID)
	if err != nil {
		Error(w, "not_found", err.Error(), http.StatusNotFound)
		return
	}

	JSON(w, summary, http.StatusOK)
}
