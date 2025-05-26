-- +goose Up
CREATE TABLE likes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,

    user_id UUID NOT NULL REFERENCES users(id),
    tweet_id UUID NOT NULL REFERENCES tweets(id),

    -- Business rule: A user can only like the same tweet once
    CONSTRAINT unique_user_tweet_like UNIQUE (user_id, tweet_id)
);

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