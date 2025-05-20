package db

import (
	"context"
	"database/sql"
	"time"
)

func InsertRefreshToken(db *sql.DB, token RefreshToken) error {
	_, err := db.ExecContext(context.Background(), `
		INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, token.UserID, token.TokenHash, token.UserAgent, token.IPAddress, token.ExpiresAt)
	return err
}

func GetRefreshTokenByHash(db *sql.DB, hash string) (*RefreshToken, error) {
	row := db.QueryRow(`
		SELECT id, user_id, token_hash, user_agent, ip_address, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`, hash)

	var rt RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.UserAgent, &rt.IPAddress, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func DeleteRefreshToken(db *sql.DB, hash string) error {
	_, err := db.Exec(`DELETE FROM refresh_tokens WHERE token_hash = $1`, hash)
	return err
}
