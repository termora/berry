package db

import (
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

	ctx, cancel := db.Context()
	defer cancel()

	err = db.Pool.QueryRow(ctx, "select token, expires from public.admin_tokens where user_id = $1", userID).Scan(&token, &expires)
	if err == pgx.ErrNoRows {
		token = GenerateToken()

		ctx, cancel := db.Context()
		defer cancel()

		commandTag, err := db.Pool.Exec(ctx, "insert into public.admin_tokens (user_id, token) values ($1, $2)", userID, token)
		if err != nil {
			return token, err
		}
		if commandTag.RowsAffected() != 1 {
			return token, ErrorNoRowsAffected
		}
		return token, err
	}
	if err != nil {
		return token, err
	}

	if expires.Before(time.Now().UTC()) {
		return token, ErrTokenExpired
	}

	return token, nil
}

// ResetToken ...
func (db *Db) ResetToken(userID string) (token string, err error) {
	ctx, cancel := db.Context()
	defer cancel()

	err = db.Pool.QueryRow(ctx, "insert into public.admin_tokens (user_id, token) values ($1, $2) on conflict (user_id) do update set token = $2, expires = (now() + interval '30 days')::timestamp returning token", userID, GenerateToken()).Scan(&token)
	return token, err
}

// ValidateToken checks if a token is valid and not expired
func (db *Db) ValidateToken(token string) (t bool) {
	ctx, cancel := db.Context()
	defer cancel()

	db.Pool.QueryRow(ctx, "select exists (select user_id from admin_tokens where token = $1 and expires > (current_timestamp at time zone 'utc') and user_id = any(select user_id from admins))", token).Scan(&t)
	return t
}
