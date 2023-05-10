package repository

import (
	"context"
	"database/sql"
	"forum/internal/models"
	"time"
)

type TokenRepository interface {
	New(userID int64, ttl time.Duration, scope string) (*models.Token, error)
	Insert(token *models.Token) error
	DeleteAllForUser(scope string, userID int64) error
}

type tokenRepo struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepo{
		db: db,
	}
}

func (r tokenRepo) New(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	token, err := models.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = r.Insert(token)
	return token, err
}

func (r tokenRepo) Insert(token *models.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r tokenRepo) DeleteAllForUser(scope string, userID int64) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, scope, userID)
	return err
}
