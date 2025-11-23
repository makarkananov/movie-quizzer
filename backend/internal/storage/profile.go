package storage

import "movie-quizzer/backend/internal/service"

func (s SQL) GetProfile(userID int64) (service.Profile, error) {
	var p service.Profile

	err := s.db.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM sessions WHERE user_id = $1) AS total_sessions,
			(SELECT COUNT(*) FROM answers WHERE session_id IN (SELECT id FROM sessions WHERE user_id = $1)) AS total_answers,
			(SELECT COUNT(*) FROM answers WHERE session_id IN (SELECT id FROM sessions WHERE user_id = $1) AND correct = TRUE) AS correct_answers,
			(SELECT COALESCE(SUM(score_delta),0) FROM answers WHERE session_id IN (SELECT id FROM sessions WHERE user_id = $1)) AS total_score
	`, userID).Scan(
		&p.TotalSessions,
		&p.TotalAnswers,
		&p.CorrectAnswers,
		&p.TotalScore,
	)
	if err != nil {
		return p, err
	}

	if p.TotalAnswers > 0 {
		p.AccuracyPercent = float64(p.CorrectAnswers) * 100 / float64(p.TotalAnswers)
	}

	return p, nil
}
