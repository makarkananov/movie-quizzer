package service

import "io"

type User struct {
	ID       int64
	Email    string
	Nickname string
}

type Session struct {
	ID             int64
	Mode           string
	Status         string
	TotalQuestions int
}

type Question struct {
	ID            int64
	Type          string
	Text          *string
	ImageURL      *string
	VideoURL      *string
	Options       []string
	CorrectAnswer string
}

type SubmitAnswer struct {
	QuestionID int64
	Answer     string
	ElapsedMs  int64
}

type AnswerResult struct {
	Correct       bool
	Score         int
	NextQuestion  *Question
	SessionStatus string
}

type SessionSummary struct {
	SessionID      int64
	TotalQuestions int
	CorrectAnswers int
	Score          int
}

type Profile struct {
	TotalSessions   int
	TotalAnswers    int
	CorrectAnswers  int
	TotalScore      int
	AccuracyPercent float64
}

type Achievement struct {
	ID          int64
	Code        string
	Title       string
	Description string
	Earned      bool
}

type LeaderboardEntry struct {
	Position int
	UserID   int64
	Nickname string
	Score    int
	Accuracy float64
}

type MediaStream struct {
	Reader      io.Reader
	ContentType string
}
