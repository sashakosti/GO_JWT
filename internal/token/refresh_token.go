package models

import "time"

type RefreshToken struct {
	ID         int       `db:"id"`
	UserID     string    `db:"user_id"`
	TokenHash  string    `db:"token_hash"`
	UserAgent  *string   `db:"user_agent"`  // nullable
	IPAddress  *string   `db:"ip_address"`  // nullable
	ExpiresAt  time.Time `db:"expires_at"`
	CreatedAt  time.Time `db:"created_at"`
}

/*func (r *RefreshToken) TableName() string {
	return "refresh_tokens"
}*/
