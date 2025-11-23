package storage

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"

	"movie-quizzer/backend/internal/service"
)

// shuffleOptions перемешивает массив вариантов ответов
// Используем Fisher-Yates shuffle для равномерного распределения
func shuffleOptions(options []string) []string {
	if len(options) <= 1 {
		return options
	}

	shuffled := make([]string, len(options))
	copy(shuffled, options)

	// Используем текущее время + дополнительную случайность для seed
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	// Fisher-Yates shuffle - гарантирует равномерное распределение
	for i := len(shuffled) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

func (s SQL) StartSession(userID int64, mode string) (service.Session, service.Question, error) {
	var sess service.Session

	// Сначала считаем, сколько вопросов доступно для данного режима
	var availableCount int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM questions WHERE type = $1
	`, mode).Scan(&availableCount)
	if err != nil {
		return sess, service.Question{}, err
	}

	// Используем минимум из доступных вопросов и 10 (или 5, если вопросов меньше)
	totalQuestions := availableCount
	if totalQuestions > 10 {
		totalQuestions = 10
	}
	if totalQuestions == 0 {
		return sess, service.Question{}, fmt.Errorf("no questions available for mode %s", mode)
	}

	// Создаем сессию
	err = s.db.QueryRow(`
		INSERT INTO sessions (user_id, mode, total_questions)
		VALUES ($1, $2, $3)
		RETURNING id, mode, status, total_questions
	`, userID, mode, totalQuestions).Scan(
		&sess.ID, &sess.Mode, &sess.Status, &sess.TotalQuestions,
	)
	if err != nil {
		return sess, service.Question{}, err
	}

	// Выбираем случайные вопросы для сессии (столько, сколько доступно, но не больше totalQuestions)
	rows, err := s.db.Query(`
		SELECT id, type, text, image_url, video_url, options, correct_answer
		FROM questions
		WHERE type = $1
		ORDER BY random()
		LIMIT $2
	`, mode, totalQuestions)
	if err != nil {
		return sess, service.Question{}, err
	}
	defer rows.Close()

	var questions []service.Question
	seenIDs := make(map[int64]bool) // Отслеживаем уже добавленные вопросы
	position := 1
	for rows.Next() {
		var q service.Question
		err := rows.Scan(
			&q.ID, &q.Type, &q.Text, &q.ImageURL, &q.VideoURL, pq.Array(&q.Options), &q.CorrectAnswer,
		)
		if err != nil {
			return sess, service.Question{}, err
		}

		// Пропускаем дубликаты (на всякий случай)
		if seenIDs[q.ID] {
			continue
		}
		seenIDs[q.ID] = true

		// Перемешиваем варианты ответов ПЕРЕД сохранением в массив
		// Это гарантирует, что правильный ответ не всегда будет на одном месте
		q.Options = shuffleOptions(q.Options)
		questions = append(questions, q)

		// Сохраняем вопрос в сессию
		_, err = s.db.Exec(`
			INSERT INTO session_questions (session_id, question_id, position)
			VALUES ($1, $2, $3)
		`, sess.ID, q.ID, position)
		if err != nil {
			return sess, service.Question{}, err
		}
		position++
	}

	if len(questions) == 0 {
		return sess, service.Question{}, fmt.Errorf("no questions available for mode %s", mode)
	}

	// Возвращаем первый вопрос
	return sess, questions[0], nil
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
	// Начисляем очки: базовые 10 за правильный ответ, бонус за скорость
	score := 0
	if isCorrect {
		score = 10
		// Бонус за скорость: если ответ дан быстрее 30 секунд, добавляем до 5 очков
		if req.ElapsedMs < 30000 {
			bonus := int(5 * (30000 - req.ElapsedMs) / 30000)
			score += bonus
		}
	}

	_, err = s.db.Exec(`
		INSERT INTO answers (session_id, question_id, user_answer, correct, score_delta, elapsed_ms)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, sessionID, req.QuestionID, req.Answer, isCorrect, score, req.ElapsedMs)
	if err != nil {
		return service.AnswerResult{}, err
	}

	// Проверяем достижение "Скоростной демон" (ответ быстрее 10 секунд)
	if isCorrect && req.ElapsedMs < 10000 {
		s.awardAchievement(userID, "speed_demon")
	}

	next, err := s.GetNextQuestion(userID, sessionID)
	if err != nil {
		// Сессия завершена - проверяем достижения
		_, _ = s.db.Exec(`
			UPDATE sessions
			SET status = 'finished', finished_at = NOW()
			WHERE id = $1
		`, sessionID)

		// Проверяем достижения после завершения сессии
		s.checkSessionAchievements(userID, sessionID)

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

	// Получаем следующий вопрос, на который еще не был дан ответ
	err := s.db.QueryRow(`
		SELECT q.id, q.type, q.text, q.image_url, q.video_url, q.options, q.correct_answer
		FROM session_questions sq
		JOIN questions q ON q.id = sq.question_id
		WHERE sq.session_id = $1
		  AND sq.question_id NOT IN (
		      SELECT question_id FROM answers WHERE session_id = $1
		  )
		ORDER BY sq.position
		LIMIT 1
	`, sessionID).Scan(
		&q.ID, &q.Type, &q.Text, &q.ImageURL, &q.VideoURL, pq.Array(&q.Options), &q.CorrectAnswer,
	)

	if err != nil {
		return q, fmt.Errorf("no more questions")
	}

	// Перемешиваем варианты ответов
	q.Options = shuffleOptions(q.Options)

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

// awardAchievement присваивает достижение пользователю, если оно еще не получено
func (s SQL) awardAchievement(userID int64, achievementCode string) {
	var achievementID int64
	err := s.db.QueryRow(`
		SELECT id FROM achievements WHERE code = $1
	`, achievementCode).Scan(&achievementID)
	if err != nil {
		return // Достижение не найдено
	}

	// Пытаемся добавить достижение (игнорируем ошибку, если уже есть)
	_, _ = s.db.Exec(`
		INSERT INTO user_achievements (user_id, achievement_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, achievement_id) DO NOTHING
	`, userID, achievementID)
}

// checkSessionAchievements проверяет и присваивает достижения после завершения сессии
func (s SQL) checkSessionAchievements(userID, sessionID int64) {
	// Проверяем, что сессия завершена с 10 вопросами
	var sessionTotalQuestions int
	err := s.db.QueryRow(`
		SELECT total_questions FROM sessions WHERE id = $1 AND status = 'finished'
	`, sessionID).Scan(&sessionTotalQuestions)
	if err != nil || sessionTotalQuestions != 10 {
		// Сессия не завершена или не имеет 10 вопросов - не засчитываем достижения
		return
	}

	// 1. Первая игра - проверяем, есть ли у пользователя другие завершенные сессии с 10 вопросами
	var otherSessions int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM sessions 
		WHERE user_id = $1 AND status = 'finished' AND total_questions = 10 AND id != $2
	`, userID, sessionID).Scan(&otherSessions)
	if err == nil && otherSessions == 0 {
		s.awardAchievement(userID, "first_game")
	}

	// 2. Идеальный раунд - все ответы правильные
	var totalQuestions, correctAnswers int
	err = s.db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM session_questions WHERE session_id = $1) AS total,
			(SELECT COUNT(*) FROM answers WHERE session_id = $1 AND correct = TRUE) AS correct
	`, sessionID).Scan(&totalQuestions, &correctAnswers)
	if err == nil && totalQuestions > 0 && correctAnswers == totalQuestions {
		s.awardAchievement(userID, "perfect_round")
	}

	// 3. Сотня - набрано 100+ очков за раунд
	var totalScore int
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(score_delta), 0) FROM answers WHERE session_id = $1
	`, sessionID).Scan(&totalScore)
	if err == nil && totalScore >= 100 {
		s.awardAchievement(userID, "century")
	}

	// 4. Ветеран - сыграно 10+ раундов с 10 вопросами
	var totalSessions int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM sessions 
		WHERE user_id = $1 AND status = 'finished' AND total_questions = 10
	`, userID).Scan(&totalSessions)
	if err == nil && totalSessions >= 10 {
		s.awardAchievement(userID, "veteran")
	}
}
