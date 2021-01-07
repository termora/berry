package db

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
)

// Errors regarding tokens
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token has expired")
)

// GenerateToken generates a 64-character token
func GenerateToken() string {
	b := make([]byte, 48)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(b)
}

// GetOrCreateToken gets or creates a token for the given user
func (db *Db) GetOrCreateToken(userID string) (token string, err error) {
	var expires time.Time
	err = db.Pool.QueryRow(context.Background(), "select token, expires from public.admin_tokens where user_id = $1").Scan(&token, &expires)
	if err == pgx.ErrNoRows {
		token = GenerateToken()
		commandTag, err := db.Pool.Exec(context.Background(), "insert into public.admin_tokens (user_id, token) values ($1, $2)", userID, token)
		if err != nil {
			return token, err
		}
		if commandTag.RowsAffected() != 1 {
			return token, ErrorNoRowsAffected
		}
		return token, err
	}

	if expires.Before(time.Now().UTC()) {
		return token, ErrTokenExpired
	}

	return token, nil
}

// ValidateToken checks if a token is valid and not expired
func (db *Db) ValidateToken(token string) (userID string, err error) {
	var expires time.Time
	err = db.Pool.QueryRow(context.Background(), "select user_id, expires from public.admin_tokens where token = $1").Scan(&userID, &expires)
	if err == pgx.ErrNoRows {
		return "", ErrInvalidToken
	}
	if err != nil {
		return "", err
	}

	if expires.Before(time.Now().UTC()) {
		return userID, ErrTokenExpired
	}

	return userID, nil
}
