-- +goose Up
CREATE TABLE tweets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    post VARCHAR(280) NOT NULL,

    creator_id UUID NOT NULL REFERENCES users(id)
);

-- +goose Down
DROP TABLE tweets;
