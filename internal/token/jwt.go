package token

import (
	"crypto/sha512"
	"encoding/hex"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // задай в .env

func GenerateAccessToken(userID string) (string, error) {
	hash := sha512.Sum512([]byte(userID))
	hashedID := hex.EncodeToString(hash[:])

	claims := jwt.MapClaims{
		"user_id": hashedID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}
