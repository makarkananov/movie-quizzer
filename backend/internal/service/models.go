package service

import "io"

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type Session struct {
	ID             int64  `json:"id"`
	Mode           string `json:"mode"`
	Status         string `json:"status"`
	TotalQuestions int    `json:"total_questions"`
}

type Question struct {
	ID            int64    `json:"id"`
	Type          string   `json:"type"`
	Text          *string  `json:"text,omitempty"`
	ImageURL      *string  `json:"image_url,omitempty"`
	VideoURL      *string  `json:"video_url,omitempty"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correct_answer"`
}

type SubmitAnswer struct {
	QuestionID int64
	Answer     string
	ElapsedMs  int64
}

type AnswerResult struct {
	Correct       bool      `json:"correct"`
	Score         int       `json:"score"`
	NextQuestion  *Question `json:"next_question,omitempty"`
	SessionStatus string    `json:"session_status"`
}

type SessionSummary struct {
	SessionID      int64 `json:"session_id"`
	TotalQuestions int   `json:"total_questions"`
	CorrectAnswers int   `json:"correct_answers"`
	Score          int   `json:"score"`
}

type Profile struct {
	TotalSessions   int     `json:"total_sessions"`
	TotalAnswers    int     `json:"total_answers"`
	CorrectAnswers  int     `json:"correct_answers"`
	TotalScore      int     `json:"total_score"`
	AccuracyPercent float64 `json:"accuracy_percent"`
}

type Achievement struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Earned      bool   `json:"earned"`
}

type LeaderboardEntry struct {
	Position int     `json:"position"`
	UserID   int64   `json:"user_id"`
	Nickname string  `json:"nickname"`
	Score    int     `json:"score"`
	Accuracy float64 `json:"accuracy"`
}

type MediaStream struct {
	Reader      io.Reader
	ContentType string
}
