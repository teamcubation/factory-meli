package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFollowSave_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test data
	followerID := uuid.New()
	followedID := uuid.New()
	followID := uuid.New()
	now := time.Now()

	// Set up expectations
	mock.ExpectQuery(`INSERT INTO follows \(follower_id, followed_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, follower_id, followed_id;`).
		WithArgs(followerID, followedID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "follower_id", "followed_id"}).
			AddRow(followID, now, now, followerID, followedID))

	// Execute the method
	ctx := context.Background()
	params := FollowSaveParams{FollowerID: followerID, FollowedID: followedID}
	follow, err := repo.Save(ctx, params)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, followID, follow.ID)
	assert.Equal(t, followerID, follow.FollowerID)
	assert.Equal(t, followedID, follow.FollowedID)
	assert.WithinDuration(t, now, follow.CreatedAt, time.Second)
	assert.WithinDuration(t, now, follow.UpdatedAt, time.Second)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowSave_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test
	followerID := uuid.New()
	followedID := uuid.New()

	// Set up expectation for database error
	mock.ExpectQuery(`INSERT INTO follows \(follower_id, followed_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, follower_id, followed_id;`).
		WithArgs(followerID, followedID).
		WillReturnError(sql.ErrConnDone) // Simulate connection error

	ctx := context.Background()
	params := FollowSaveParams{FollowerID: followerID, FollowedID: followedID}
	follow, err := repo.Save(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, follow)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowSave_DuplicateFollowError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test data
	followerID := uuid.New()
	followedID := uuid.New()

	// Simulate a unique constraint violation error
	expectedErrorMsg := `pq: duplicate key value violates unique constraint "follows_follower_id_followed_id_key"`
	mock.ExpectQuery(`INSERT INTO follows \(follower_id, followed_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, follower_id, followed_id;`).
		WithArgs(followerID, followedID).
		WillReturnError(errors.New(expectedErrorMsg))

	ctx := context.Background()
	params := FollowSaveParams{FollowerID: followerID, FollowedID: followedID}
	follow, err := repo.Save(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, follow)
	assert.Contains(t, err.Error(), "duplicate key value violates unique constraint")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowSave_ForeignKeyViolation(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test data
	followerID := uuid.New() // Non-existent user
	followedID := uuid.New() // Non-existent user

	// Simulate a foreign key violation error
	expectedErrorMsg := `pq: insert or update on table "follows" violates foreign key constraint "follows_follower_id_fkey"`
	mock.ExpectQuery(`INSERT INTO follows \(follower_id, followed_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, follower_id, followed_id;`).
		WithArgs(followerID, followedID).
		WillReturnError(errors.New(expectedErrorMsg))

	ctx := context.Background()
	params := FollowSaveParams{FollowerID: followerID, FollowedID: followedID}
	follow, err := repo.Save(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, follow)
	assert.Contains(t, err.Error(), "violates foreign key constraint")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowDelete_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test data
	followerID := uuid.New()
	followedID := uuid.New()

	// Set up expectation for successful update (1 row affected)
	mock.ExpectExec(`UPDATE follows SET deleted_at = NOW\(\) WHERE follower_id = \$1 AND followed_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(followerID, followedID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	ctx := context.Background()
	params := FollowDeleteParams{FollowerID: followerID, FollowedID: followedID}
	err := repo.Delete(ctx, params)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowDelete_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test data
	followerID := uuid.New()
	followedID := uuid.New()

	// Set up expectation for no rows affected
	mock.ExpectExec(`UPDATE follows SET deleted_at = NOW\(\) WHERE follower_id = \$1 AND followed_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(followerID, followedID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	ctx := context.Background()
	params := FollowDeleteParams{FollowerID: followerID, FollowedID: followedID}
	err := repo.Delete(ctx, params)

	assert.Error(t, err)
	assert.EqualError(t, err, "follow relationship not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowDelete_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test data
	followerID := uuid.New()
	followedID := uuid.New()

	// Set up expectation for database error during ExecContext
	mock.ExpectExec(`UPDATE follows SET deleted_at = NOW\(\) WHERE follower_id = \$1 AND followed_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(followerID, followedID).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := FollowDeleteParams{FollowerID: followerID, FollowedID: followedID}
	err := repo.Delete(ctx, params)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowDelete_RowsAffectedError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test data
	followerID := uuid.New()
	followedID := uuid.New()

	// Set up expectation for an error when calling RowsAffected
	mock.ExpectExec(`UPDATE follows SET deleted_at = NOW\(\) WHERE follower_id = \$1 AND followed_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(followerID, followedID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	ctx := context.Background()
	params := FollowDeleteParams{FollowerID: followerID, FollowedID: followedID}
	err := repo.Delete(ctx, params)

	assert.Error(t, err)
	assert.EqualError(t, err, "rows affected error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowFindByIds_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	// Test data
	followerID := uuid.New()
	followedID := uuid.New()
	followID := uuid.New()
	now := time.Now()

	// Set up expectations
	mock.ExpectQuery(`SELECT id, created_at, updated_at, follower_id, followed_id FROM follows WHERE follower_id = \$1 AND followed_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(followerID, followedID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "follower_id", "followed_id"}).
			AddRow(followID, now, now, followerID, followedID))

	// Execute the method
	ctx := context.Background()
	params := FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}
	follow, err := repo.FindByIds(ctx, params)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, followID, follow.ID)
	assert.Equal(t, followerID, follow.FollowerID)
	assert.Equal(t, followedID, follow.FollowedID)
	assert.WithinDuration(t, now, follow.CreatedAt, time.Second)
	assert.WithinDuration(t, now, follow.UpdatedAt, time.Second)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowFindByIds_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	followerID := uuid.New()
	followedID := uuid.New()

	// Set up expectation for no rows found
	mock.ExpectQuery(`SELECT id, created_at, updated_at, follower_id, followed_id FROM follows WHERE follower_id = \$1 AND followed_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(followerID, followedID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	params := FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}
	follow, err := repo.FindByIds(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, follow)
	assert.Equal(t, sql.ErrNoRows, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFollowFindByIds_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) FollowRepository {
		return NewPostgresFollowRepository(db)
	})
	defer db.Close()

	followerID := uuid.New()
	followedID := uuid.New()

	// Set up expectation for database connection error
	mock.ExpectQuery(`SELECT id, created_at, updated_at, follower_id, followed_id FROM follows WHERE follower_id = \$1 AND followed_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(followerID, followedID).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := FollowFindByIdsParams{FollowerID: followerID, FollowedID: followedID}
	follow, err := repo.FindByIds(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, follow)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
