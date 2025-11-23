package storage

import "movie-quizzer/backend/internal/service"

func (s SQL) GetAchievements(userID int64) ([]service.Achievement, error) {
	rows, err := s.db.Query(`
		SELECT a.id, a.code, a.title, a.description,
		       ua.user_id IS NOT NULL AS earned
		FROM achievements a
		LEFT JOIN user_achievements ua 
		       ON ua.achievement_id = a.id AND ua.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []service.Achievement
	for rows.Next() {
		var a service.Achievement
		err := rows.Scan(&a.ID, &a.Code, &a.Title, &a.Description, &a.Earned)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}

	return out, nil
}

func (s SQL) GetGlobalLeaderboard(limit int) ([]service.LeaderboardEntry, error) {
	rows, err := s.db.Query(`
		SELECT 
		    u.id,
		    u.nickname,
		    COALESCE(SUM(a.score_delta), 0) AS score,
		    CASE
		        WHEN COUNT(a.id) = 0 THEN 0
		        ELSE SUM(CASE WHEN a.correct = TRUE THEN 1 ELSE 0 END)::float / COUNT(a.id) * 100
		    END AS accuracy
		FROM users u
		LEFT JOIN sessions s ON s.user_id = u.id AND s.status = 'finished' AND s.total_questions = 10
		LEFT JOIN answers a ON a.session_id = s.id
		GROUP BY u.id, u.nickname
		ORDER BY score DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []service.LeaderboardEntry
	i := 1

	for rows.Next() {
		var e service.LeaderboardEntry
		err := rows.Scan(&e.UserID, &e.Nickname, &e.Score, &e.Accuracy)
		if err != nil {
			return nil, err
		}
		e.Position = i
		i++
		out = append(out, e)
	}

	return out, nil
}

func (s SQL) GetLeaderboardEntry(userID int64) (service.LeaderboardEntry, error) {
	var e service.LeaderboardEntry

	err := s.db.QueryRow(`
		WITH scores AS (
			SELECT 
				u.id,
				u.nickname,
				COALESCE(SUM(a.score_delta), 0) AS score,
				CASE
					WHEN COUNT(a.id) = 0 THEN 0
					ELSE SUM(CASE WHEN a.correct = TRUE THEN 1 ELSE 0 END)::float / COUNT(a.id) * 100
				END AS accuracy
			FROM users u
			LEFT JOIN sessions s ON s.user_id = u.id AND s.status = 'finished' AND s.total_questions = 10
			LEFT JOIN answers a ON a.session_id = s.id
			GROUP BY u.id, u.nickname
		),
		ranked AS (
			SELECT *,
			       RANK() OVER (ORDER BY score DESC) AS position
			FROM scores
		)
		SELECT position, id, nickname, score, accuracy
		FROM ranked
		WHERE id = $1
	`, userID).Scan(&e.Position, &e.UserID, &e.Nickname, &e.Score, &e.Accuracy)

	return e, err
}
