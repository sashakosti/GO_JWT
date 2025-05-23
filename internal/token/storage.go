package token

import (
	"database/sql"
	"errors"
)

// Storage defines the interface for token storage
type Storage interface {
	Store(token *RefreshToken) error
	GetByHash(hash string) (*RefreshToken, error)
	Delete(id string) error
	DeleteExpired() error
}

// PostgresStorage implements Storage for PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) Store(token *RefreshToken) error {
	// TODO: Implement PostgreSQL storage
	return errors.New("not implemented")
}

func (s *PostgresStorage) GetByHash(hash string) (*RefreshToken, error) {
	// TODO: Implement PostgreSQL retrieval
	return nil, errors.New("not implemented")
}

func (s *PostgresStorage) Delete(id string) error {
	// TODO: Implement PostgreSQL deletion
	return errors.New("not implemented")
}

func (s *PostgresStorage) DeleteExpired() error {
	// TODO: Implement PostgreSQL cleanup
	return errors.New("not implemented")
}
