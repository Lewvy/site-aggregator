-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1;

-- name: DropRows :exec
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

-- name: ListUsers :many
SELECT name FROM users;

-- name: InsertFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListFeeds :many
SELECT feeds.name, feeds.url, users.name 
FROM feeds 
INNER JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
WITH ins AS (
    INSERT INTO feed_follows (user_id, feed_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, feed_id) DO NOTHING
    RETURNING id, feed_follows.user_id, feed_follows.feed_id, feed_follows.created_at, feed_follows.updated_at
)
SELECT
    ff.id,
    ff.user_id,
    ff.feed_id,
    ff.created_at,
    ff.updated_at,
    u.name AS user_name,
    f.name AS feed_name
FROM
    (SELECT * FROM ins
     UNION ALL
     SELECT id, user_id, feed_id, created_at, updated_at
     FROM feed_follows WHERE feed_follows.user_id = $1 AND feed_follows.feed_id = $2) ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id
LIMIT 1;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.id,
    ff.user_id,
    ff.feed_id,
    ff.created_at,
    ff.updated_at,
    u.name AS user_name,
    f.name AS feed_name
FROM
    feed_follows ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id
WHERE
    ff.user_id = $1;

-- name: GetFeedNamesUserIsFollowing :many
SELECT
    f.name AS feed_name
FROM
    feed_follows ff
JOIN feeds f ON ff.feed_id = f.id
WHERE
    ff.user_id = $1;

-- name: AddFeedAndFollow :one
WITH ins_feed AS (
    INSERT INTO feeds (url, name)
    VALUES ($1, $2)
    ON CONFLICT (url) DO UPDATE SET name = EXCLUDED.name
    RETURNING id, name
),
sel_feed AS (
    SELECT id, name FROM feeds WHERE url = $1
    UNION ALL
    SELECT id, name FROM ins_feed
    LIMIT 1
),
ins_follow AS (
    INSERT INTO feed_follows (user_id, feed_id)
    SELECT $3, id FROM sel_feed
    ON CONFLICT (user_id, feed_id) DO NOTHING
    RETURNING id
)
SELECT
    f.name AS feed_name,
    u.name AS user_name
FROM
    sel_feed f,
    users u
WHERE
    u.id = $3
LIMIT 1;
