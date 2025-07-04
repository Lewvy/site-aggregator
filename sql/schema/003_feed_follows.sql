-- +goose Up
CREATE TABLE feed_follows (
    id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id INTEGER NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_feed UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
