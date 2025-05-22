package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// RepositoryConstructor defines a function type for creating repositories
type RepositoryConstructor[T any] func(*sql.DB) T

// setupMockDB creates a mock database and repository of the specified type
func setupMockDB[T any](t *testing.T, constructor RepositoryConstructor[T]) (*sql.DB, sqlmock.Sqlmock, T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	
	repo := constructor(db)
	return db, mock, repo
}
