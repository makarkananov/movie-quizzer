package storage

import (
	"fmt"

	"movie-quizzer/backend/internal/service"
)

func (s SQL) CreateUser(email, password, nickname string) error {
	_, err := s.db.Exec(`
		INSERT INTO users (email, password_hash, nickname)
		VALUES ($1, $2, $3)
	`, email, password, nickname)
	return err
}

func (s SQL) LoginUser(email, password string) (string, error) {
	var id int64
	var hash string

	err := s.db.QueryRow(`
		SELECT id, password_hash FROM users WHERE email = $1
	`, email).Scan(&id, &hash)
	if err != nil {
		return "", err
	}

	if hash != password {
		return "", fmt.Errorf("wrong_credentials")
	}

	return fmt.Sprintf("token-%d", id), nil
}

func (s SQL) GetUserFromToken(token string) (service.User, error) {
	var id int64
	_, err := fmt.Sscanf(token, "token-%d", &id)
	if err != nil {
		return service.User{}, fmt.Errorf("invalid token")
	}

	var u service.User
	err = s.db.QueryRow(`
		SELECT id, email, nickname
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.Nickname)

	return u, err
}
