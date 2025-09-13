package postgre

import (
	"context"
	"fmt"
	"idea-store-auth/internal/domain/models"
	"log/slog"
)

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	slog.Info("storage started to SaveUser")

	const query = `
		INSERT INTO users(email, pass_hash) 
		VALUES($1,$2)
		RETURNING id;
	`

	var lastInsertId int64
	err := s.db.QueryRow(ctx, query, email, passHash).Scan(&lastInsertId)
	if err != nil {
		slog.Error("storage SaveUser error: " + err.Error())
		return emptyValue, fmt.Errorf("storage SaveUser error: %v", err.Error())
	}

	return lastInsertId, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	slog.Info("storage started to SaveUser")

	const query = `
		SELECT id,email,pass_hash 
		FROM users 
		WHERE email = $1;
	`

	var user models.User
	err := s.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		slog.Error("storage User error: " + err.Error())
		return models.User{}, fmt.Errorf("storage User error: %v", err.Error())
	}

	return user, nil
}
