-- +goose Up
CREATE TABLE feeds (
  id serial PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL UNIQUE,
  user_id uuid not null,
  FOREIGN KEY(user_id)
       REFERENCES users(id)
       ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;

