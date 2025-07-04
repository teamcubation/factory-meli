-- +goose Up
CREATE TABLE tweets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    post VARCHAR(280) NOT NULL,

    creator_id UUID NOT NULL REFERENCES users(id)
);

-- Critical composite index for timeline query: covers creator_id filtering + ordering by created_at
CREATE INDEX idx_tweets_creator_created_active ON tweets (creator_id, created_at DESC) WHERE deleted_at IS NULL;

-- Additional index for general creator queries
CREATE INDEX idx_tweets_creator_id ON tweets (creator_id) WHERE deleted_at IS NULL;

-- Index for soft delete filtering
CREATE INDEX idx_tweets_deleted_at ON tweets (deleted_at) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE tweets;
