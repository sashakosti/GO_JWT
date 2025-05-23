package token

import (
	/*"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"*/
	"time"
)

type RefreshTokenManager struct {
	storage    Storage
	expiration time.Duration
}

func NewRefreshTokenManager(storage Storage, expiration time.Duration) *RefreshTokenManager {
	return &RefreshTokenManager{
		storage:    storage,
		expiration: expiration,
	}
}

func (m *RefreshTokenManager) Generate(userID string) (string, error) {
	token, err := generateRandomString(32)
	if err != nil {
		return "", err
	}

	hashedToken, err := hashToken(token)
	if err != nil {
		return "", err
	}

	refreshToken := &RefreshToken{
		UserID:    userID,
		TokenHash: hashedToken,
		ExpiresAt: time.Now().Add(m.expiration),
		CreatedAt: time.Now(),
	}

	if err := m.storage.Store(refreshToken); err != nil {
		return "", err
	}

	return token, nil
}

func (m *RefreshTokenManager) Validate(userID, token string) error {
	hashedToken, err := hashToken(token)
	if err != nil {
		return err
	}

	storedToken, err := m.storage.GetByHash(hashedToken)
	if err != nil {
		return err
	}

	if storedToken.UserID != userID {
		return ErrInvalidToken
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return ErrExpiredToken
	}

	return m.storage.Delete(storedToken.ID)
}
