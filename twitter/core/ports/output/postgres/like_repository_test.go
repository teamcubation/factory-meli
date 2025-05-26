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

func TestLikeHas_Success_Exists(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM likes WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL\);`).
		WithArgs(userID, tweetID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	ctx := context.Background()
	params := LikeHasParams{UserID: userID, TweetID: tweetID}
	exists, err := repo.HasLike(ctx, params)

	require.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikeHas_Success_NotExists(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM likes WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL\);`).
		WithArgs(userID, tweetID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	ctx := context.Background()
	params := LikeHasParams{UserID: userID, TweetID: tweetID}
	exists, err := repo.HasLike(ctx, params)

	require.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikeHas_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM likes WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL\);`).
		WithArgs(userID, tweetID).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := LikeHasParams{UserID: userID, TweetID: tweetID}
	exists, err := repo.HasLike(ctx, params)

	assert.Error(t, err)
	assert.False(t, exists)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddLike_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()
	likeID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`INSERT INTO likes \(user_id, tweet_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, user_id, tweet_id;`).
		WithArgs(userID, tweetID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "tweet_id"}).
			AddRow(likeID, now, now, userID, tweetID))

	ctx := context.Background()
	params := LikeAddParams{UserID: userID, TweetID: tweetID}
	like, err := repo.AddLike(ctx, params)

	require.NoError(t, err)
	assert.Equal(t, likeID, like.ID)
	assert.Equal(t, userID, like.UserID)
	assert.Equal(t, tweetID, like.TweetID)
	assert.WithinDuration(t, now, like.CreatedAt, time.Second)
	assert.WithinDuration(t, now, like.UpdatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddLike_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectQuery(`INSERT INTO likes \(user_id, tweet_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, user_id, tweet_id;`).
		WithArgs(userID, tweetID).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := LikeAddParams{UserID: userID, TweetID: tweetID}
	like, err := repo.AddLike(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, like)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddLike_ForeignKeyViolation(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	expectedErrorMsg := `pq: insert or update on table "likes" violates foreign key constraint "likes_user_id_fkey"`
	mock.ExpectQuery(`INSERT INTO likes \(user_id, tweet_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, user_id, tweet_id;`).
		WithArgs(userID, tweetID).
		WillReturnError(errors.New(expectedErrorMsg))

	ctx := context.Background()
	params := LikeAddParams{UserID: userID, TweetID: tweetID}
	like, err := repo.AddLike(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, like)
	assert.Contains(t, err.Error(), "violates foreign key constraint")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveLike_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectExec(`UPDATE likes SET deleted_at = NOW\(\), updated_at = NOW\(\) WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(userID, tweetID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	ctx := context.Background()
	params := LikeRemoveParams{UserID: userID, TweetID: tweetID}
	err := repo.RemoveLike(ctx, params)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveLike_NoRowsAffected(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectExec(`UPDATE likes SET deleted_at = NOW\(\), updated_at = NOW\(\) WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(userID, tweetID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	ctx := context.Background()
	params := LikeRemoveParams{UserID: userID, TweetID: tweetID}
	err := repo.RemoveLike(ctx, params)

	require.NoError(t, err) // No error is returned if no rows are affected
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveLike_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectExec(`UPDATE likes SET deleted_at = NOW\(\), updated_at = NOW\(\) WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(userID, tweetID).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	params := LikeRemoveParams{UserID: userID, TweetID: tweetID}
	err := repo.RemoveLike(ctx, params)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveLike_RowsAffectedError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) LikeRepository {
		return NewPostgresLikeRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectExec(`UPDATE likes SET deleted_at = NOW\(\), updated_at = NOW\(\) WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL;`).
		WithArgs(userID, tweetID).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

	ctx := context.Background()
	params := LikeRemoveParams{UserID: userID, TweetID: tweetID}
	err := repo.RemoveLike(ctx, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rows affected error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
