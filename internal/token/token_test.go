package token_test

import (
	"github.com/sashakosti/auth-service/internal/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// MockStorage implements the Storage interface for testing
type MockStorage struct {
	tokens map[string]*token.RefreshToken
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		tokens: make(map[string]*token.RefreshToken),
	}
}

func (m *MockStorage) Store(t *token.RefreshToken) error {
	m.tokens[t.TokenHash] = t
	return nil
}

func (m *MockStorage) GetByHash(hash string) (*token.RefreshToken, error) {
	if t, ok := m.tokens[hash]; ok {
		return t, nil
	}
	return nil, token.ErrTokenNotFound
}

func (m *MockStorage) Delete(id string) error {
	for hash, t := range m.tokens {
		if t.ID == id {
			delete(m.tokens, hash)
			return nil
		}
	}
	return token.ErrTokenNotFound
}

func (m *MockStorage) DeleteExpired() error {
	now := time.Now()
	for hash, t := range m.tokens {
		if t.ExpiresAt.Before(now) {
			delete(m.tokens, hash)
		}
	}
	return nil
}

func TestTokenFlow(t *testing.T) {
	// Setup
	storage := NewMockStorage()
	accessExp := 15 * time.Minute
	refreshExp := 7 * 24 * time.Hour

	manager := token.NewTokenManager(
		"test-secret-123",
		accessExp,
		refreshExp,
		storage,
	)

	userID := "user-123"

	t.Run("generate tokens", func(t *testing.T) {
		// Generate token pair
		pair, err := manager.GenerateTokens(userID)
		require.NoError(t, err)
		require.NotEmpty(t, pair.AccessToken)
		require.NotEmpty(t, pair.RefreshToken)

		// Validate access token
		claims, err := manager.ValidateAccessToken(pair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims)

		// Try to validate with wrong secret (should fail)
		wrongManager := token.NewTokenManager("wrong-secret", accessExp, refreshExp, storage)
		_, err = wrongManager.ValidateAccessToken(pair.AccessToken)
		assert.Error(t, err)
	})

	t.Run("refresh tokens", func(t *testing.T) {
		// Generate initial token pair
		pair, err := manager.GenerateTokens(userID)
		require.NoError(t, err)

		// Refresh tokens
		newPair, err := manager.RefreshTokens(userID, pair.RefreshToken)
		require.NoError(t, err)
		require.NotEqual(t, pair.AccessToken, newPair.AccessToken)
		require.NotEqual(t, pair.RefreshToken, newPair.RefreshToken)

		// Old refresh token should be invalid now
		_, err = manager.RefreshTokens(userID, pair.RefreshToken)
		assert.ErrorIs(t, err, token.ErrTokenNotFound)
	})

	t.Run("expired token", func(t *testing.T) {
		// Create a manager with very short expiration
		tempManager := token.NewTokenManager(
			"test-secret-123",
			1*time.Millisecond, // Very short access token lifetime
			1*time.Millisecond, // Very short refresh token lifetime
			storage,
		)

		// Generate and wait for token to expire
		pair, err := tempManager.GenerateTokens(userID)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		// Access token should be expired
		_, err = tempManager.ValidateAccessToken(pair.AccessToken)
		assert.ErrorIs(t, err, token.ErrExpiredToken)

		// Refresh token should also be expired
		_, err = tempManager.RefreshTokens(userID, pair.RefreshToken)
		assert.ErrorIs(t, err, token.ErrTokenNotFound)
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		_, err := manager.RefreshTokens(userID, "invalid-token")
		assert.ErrorIs(t, err, token.ErrTokenNotFound)
	})
}
