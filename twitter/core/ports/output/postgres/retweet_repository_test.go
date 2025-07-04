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

func TestRetweetHas_Success_Exists(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) RetweetRepository {
		return NewPostgresRetweetRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM retweets WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL\);`).
		WithArgs(userID, tweetID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	ctx := context.Background()
	params := RetweetHasParams{UserID: userID, TweetID: tweetID}
	exists, err := repo.HasRetweet(ctx, params)

	require.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetweetHas_Success_NotExists(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) RetweetRepository {
		return NewPostgresRetweetRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM retweets WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL\);`).
		WithArgs(userID, tweetID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	ctx := context.Background()
	params := RetweetHasParams{UserID: userID, TweetID: tweetID}
	exists, err := repo.HasRetweet(ctx, params)

	require.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetweetHas_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) RetweetRepository {
		return NewPostgresRetweetRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM retweets WHERE user_id = \$1 AND tweet_id = \$2 AND deleted_at IS NULL\);`).
		WithArgs(userID, tweetID).
		WillReturnError(sql.ErrConnDone) // Simulate a database error

	ctx := context.Background()
	params := RetweetHasParams{UserID: userID, TweetID: tweetID}
	exists, err := repo.HasRetweet(ctx, params)

	assert.Error(t, err)
	assert.False(t, exists)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddRetweet_Success(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) RetweetRepository {
		return NewPostgresRetweetRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()
	retweetID := uuid.New()
	now := time.Now()

	mock.ExpectQuery(`INSERT INTO retweets \(user_id, tweet_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, user_id, tweet_id;`).
		WithArgs(userID, tweetID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "user_id", "tweet_id"}).
			AddRow(retweetID, now, now, userID, tweetID))

	ctx := context.Background()
	params := RetweetAddParams{UserID: userID, TweetID: tweetID}
	retweet, err := repo.AddRetweet(ctx, params)

	require.NoError(t, err)
	assert.Equal(t, retweetID, retweet.ID)
	assert.Equal(t, userID, retweet.UserID)
	assert.Equal(t, tweetID, retweet.TweetID)
	assert.WithinDuration(t, now, retweet.CreatedAt, time.Second)
	assert.WithinDuration(t, now, retweet.UpdatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddRetweet_DatabaseError(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) RetweetRepository {
		return NewPostgresRetweetRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	mock.ExpectQuery(`INSERT INTO retweets \(user_id, tweet_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, user_id, tweet_id;`).
		WithArgs(userID, tweetID).
		WillReturnError(sql.ErrConnDone) // Simulate a database error

	ctx := context.Background()
	params := RetweetAddParams{UserID: userID, TweetID: tweetID}
	retweet, err := repo.AddRetweet(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, retweet)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddRetweet_ForeignKeyViolation(t *testing.T) {
	t.Parallel()

	db, mock, repo := setupMockDB(t, func(db *sql.DB) RetweetRepository {
		return NewPostgresRetweetRepository(db)
	})
	defer db.Close()

	userID := uuid.New()
	tweetID := uuid.New()

	expectedErrorMsg := `pq: insert or update on table "retweets" violates foreign key constraint "retweets_user_id_fkey"`
	mock.ExpectQuery(`INSERT INTO retweets \(user_id, tweet_id\) VALUES \(\$1, \$2\) RETURNING id, created_at, updated_at, user_id, tweet_id;`).
		WithArgs(userID, tweetID).
		WillReturnError(errors.New(expectedErrorMsg))

	ctx := context.Background()
	params := RetweetAddParams{UserID: userID, TweetID: tweetID}
	retweet, err := repo.AddRetweet(ctx, params)

	assert.Error(t, err)
	assert.Nil(t, retweet)
	assert.Contains(t, err.Error(), "violates foreign key constraint")
	assert.NoError(t, mock.ExpectationsWereMet())
}
