-- +goose Up
CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    user_id UUID NOT NULL REFERENCES users(id),
    tweet_id UUID NOT NULL REFERENCES tweets(id)
);

-- Instead of a UNIQUE constraint, we use a partial unique index that only applies to non-deleted rows
CREATE UNIQUE INDEX idx_user_tweet_like_unique ON likes (user_id, tweet_id) WHERE deleted_at IS NULL;

-- Index for finding all likes by a user
CREATE INDEX idx_likes_user_id ON likes (user_id) WHERE deleted_at IS NULL;

-- Index for finding all likes of a tweet (for counting likes)
CREATE INDEX idx_likes_tweet_id ON likes (tweet_id) WHERE deleted_at IS NULL;

-- Composite index for user timeline queries (user's liked tweets ordered by creation date)
CREATE INDEX idx_likes_user_created_active ON likes (user_id, created_at DESC) WHERE deleted_at IS NULL;

-- Index for soft delete filtering
CREATE INDEX idx_likes_deleted_at ON likes (deleted_at) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE likes;