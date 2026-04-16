-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE name = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT name
FROM users;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetUserID :one
SELECT id
FROM users
WHERE name = $1;

-- name: GetFeedID :one
SELECT id
FROM feeds
WHERE url = $1;

-- name: GetFeedsInfo :many
SELECT 
    feeds.name,
    feeds.url,
    users.name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT inserted_feed_follow.*,
feeds.name AS feed_name,
users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name AS feed_name
FROM feed_follows 
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1; 

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE user_id = $1
AND feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET 
    last_fetched_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;