package token

import "time"

type TokenManager interface {
	GenerateTokens(userID string) (*TokenPair, error)
	RefreshTokens(userID, refreshToken string) (*TokenPair, error)
	ValidateAccessToken(tokenString string) (string, error)
}

type tokenManager struct {
	access  *AccessTokenManager
	refresh *RefreshTokenManager
}

func NewTokenManager(accessSecret string, accessExp, refreshExp time.Duration, storage Storage) TokenManager {
	return &tokenManager{
		access:  NewAccessTokenManager(accessSecret, accessExp),
		refresh: NewRefreshTokenManager(storage, refreshExp),
	}
}

func (m *tokenManager) GenerateTokens(userID string) (*TokenPair, error) {
	accessToken, err := m.access.Generate(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.refresh.Generate(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(m.access.expiration),
	}, nil
}

func (m *tokenManager) RefreshTokens(userID, refreshToken string) (*TokenPair, error) {
	if err := m.refresh.Validate(userID, refreshToken); err != nil {
		return nil, err
	}
	return m.GenerateTokens(userID)
}

func (m *tokenManager) ValidateAccessToken(tokenString string) (string, error) {
	return m.access.Validate(tokenString)
}
