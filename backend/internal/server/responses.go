package server

import (
	"encoding/json"
	"net/http"
)

type errRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, code, msg string, status int) {
	JSON(w, errRes{Code: code, Message: msg}, status)
}
