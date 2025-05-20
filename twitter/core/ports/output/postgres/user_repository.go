package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/twitter-tq/vinofsteel/core/domain/models"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	Save(ctx context.Context, params UserSaveParams) (*models.User, error)
	FindByEmail(ctx context.Context, params UserFindByEmailParams) (*models.User, error)
}

// PostgresUserRepository implements UserRepository interface with PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Save inserts a new user into the database
type UserSaveParams struct {
	Email string `json:"email"`
}

func (r *PostgresUserRepository) Save(ctx context.Context, params UserSaveParams) (*models.User, error) {
	slog.InfoContext(ctx, "Saving a new user to the database", "email", params.Email)
	query := `INSERT INTO users (email) VALUES ($1) RETURNING id, created_at, updated_at, email;`

	var user models.User
	err := r.db.QueryRowContext(
		ctx,
		query,
		params.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
	)
	if err != nil {
		slog.ErrorContext(ctx, "Error saving user to db", "error", err)
		return nil, err
	}

	return &user, nil
}

// FindByEmail retrieves a user by their email
type UserFindByEmailParams struct {
	Email string `json:"email"`
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, params UserFindByEmailParams) (*models.User, error) {
	slog.InfoContext(ctx, "Finding an user by email", "email", params.Email)
	query := `SELECT id, created_at, updated_at, email FROM users WHERE email = $1 AND deleted_at IS NULL;`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, params.Email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
