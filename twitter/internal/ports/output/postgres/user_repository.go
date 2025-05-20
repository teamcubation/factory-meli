package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twitter-tq/vinofsteel/internal/domain/models"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	Save(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
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
func (r *PostgresUserRepository) Save(ctx context.Context, user *models.User) error {
	// Set the created time if not already set
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	user.UpdatedAt = time.Now()

	// If ID is not set, generate a new UUID
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	query := `
		INSERT INTO users (id, created_at, updated_at, email)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Email,
	)

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// FindByID retrieves a user by their ID
func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, email
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
		&u.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil when no user is found
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return &u, nil
}

// FindByEmail retrieves a user by their email
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, email
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
		&u.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil when no user is found
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &u, nil
}

// FindAll retrieves all non-deleted users
func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, created_at, updated_at, deleted_at, email
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.DeletedAt,
			&u.Email,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET updated_at = $1, email = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.UpdatedAt,
		user.Email,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// Delete soft-deletes a user by setting the deleted_at timestamp
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}
