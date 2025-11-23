package storage

import (
	"fmt"

	"github.com/lib/pq"

	"movie-quizzer/backend/internal/service"
)

func (s SQL) StartSession(userID int64, mode string) (service.Session, service.Question, error) {
	var sess service.Session

	err := s.db.QueryRow(`
		INSERT INTO sessions (user_id, mode)
		VALUES ($1, $2)
		RETURNING id, mode, status, total_questions
	`, userID, mode).Scan(
		&sess.ID, &sess.Mode, &sess.Status, &sess.TotalQuestions,
	)
	if err != nil {
		return sess, service.Question{}, err
	}

	var q service.Question
	err = s.db.QueryRow(`
		SELECT id, type, text, image_url, video_url, options, correct_answer
		FROM questions
		WHERE type = $1
		ORDER BY random()
		LIMIT 1
	`, mode).Scan(
		&q.ID, &q.Type, &q.Text, &q.ImageURL, &q.VideoURL, pq.Array(&q.Options), &q.CorrectAnswer,
	)
	if err != nil {
		return sess, q, err
	}

	_, err = s.db.Exec(`
		INSERT INTO session_questions (session_id, question_id, position)
		VALUES ($1, $2, 1)
	`, sess.ID, q.ID)

	return sess, q, err
}

func (s SQL) SubmitAnswer(userID, sessionID int64, req service.SubmitAnswer) (service.AnswerResult, error) {
	var correctAnswer string

	err := s.db.QueryRow(`
		SELECT correct_answer FROM questions WHERE id = $1
	`, req.QuestionID).Scan(&correctAnswer)
	if err != nil {
		return service.AnswerResult{}, err
	}

	isCorrect := req.Answer == correctAnswer
	score := 10
	if !isCorrect {
		score = 0
	}

	_, err = s.db.Exec(`
		INSERT INTO answers (session_id, question_id, user_answer, correct, score_delta, elapsed_ms)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, sessionID, req.QuestionID, req.Answer, isCorrect, score, req.ElapsedMs)
	if err != nil {
		return service.AnswerResult{}, err
	}

	next, err := s.GetNextQuestion(userID, sessionID)
	if err != nil {
		_, _ = s.db.Exec(`
			UPDATE sessions
			SET status = 'finished', finished_at = NOW()
			WHERE id = $1
		`, sessionID)

		return service.AnswerResult{
			Correct:       isCorrect,
			Score:         score,
			SessionStatus: "finished",
			NextQuestion:  nil,
		}, nil
	}

	return service.AnswerResult{
		Correct:       isCorrect,
		Score:         score,
		SessionStatus: "in_progress",
		NextQuestion:  &next,
	}, nil
}

func (s SQL) GetNextQuestion(userID, sessionID int64) (service.Question, error) {
	var q service.Question

	err := s.db.QueryRow(`
		SELECT q.id, q.type, q.text, q.image_url, q.video_url, q.options, q.correct_answer
		FROM session_questions sq
		JOIN questions q ON q.id = sq.question_id
		WHERE sq.session_id = $1
		ORDER BY sq.position
		LIMIT 1
	`, sessionID).Scan(
		&q.ID, &q.Type, &q.Text, &q.ImageURL, &q.VideoURL, pq.Array(&q.Options), &q.CorrectAnswer,
	)

	if err != nil {
		return q, fmt.Errorf("no more questions")
	}

	return q, nil
}

func (s SQL) GetSessionSummary(userID, sessionID int64) (service.SessionSummary, error) {
	var sum service.SessionSummary

	err := s.db.QueryRow(`
		SELECT 
			s.id,
			(SELECT COUNT(*) FROM session_questions WHERE session_id = s.id) AS total_questions,
			(SELECT COUNT(*) FROM answers WHERE session_id = s.id AND correct = TRUE) AS correct_answers,
			(SELECT COALESCE(SUM(score_delta),0) FROM answers WHERE session_id = s.id) AS score
		FROM sessions s
		WHERE s.id = $1
	`, sessionID).Scan(&sum.SessionID, &sum.TotalQuestions, &sum.CorrectAnswers, &sum.Score)

	return sum, err
}
