package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AccessTokenManager struct {
	secret     []byte
	expiration time.Duration
}

func NewAccessTokenManager(secret string, expiration time.Duration) *AccessTokenManager {
	return &AccessTokenManager{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (m *AccessTokenManager) Generate(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(m.secret)
}

func (m *AccessTokenManager) Validate(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", ErrInvalidToken
}
