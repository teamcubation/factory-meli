CREATE TABLE
    IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE
    );

CREATE TABLE
    IF NOT EXISTS follows (
        user_id UUID NOT NULL,
        follow_id UUID NOT NULL,
        PRIMARY KEY (user_id, follow_id),
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
        FOREIGN KEY (follow_id) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS tweets (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        content VARCHAR(280) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );