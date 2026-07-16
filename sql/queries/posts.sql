-- name: AddPostsToDatabase :many
INSERT INTO posts (id, feed_id, title, url, description, published_at, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsByFeed :many
SELECT 
    f.id AS feed_id,
    f.name AS feed_name,
    f.url AS feed_url,
    p.id AS post_id,
    p.title AS post_title,
    p.url AS post_url,
    p.description AS post_description,
    p.published_at AS post_published_at,
    p.created_at AS post_created_at,
    p.updated_at AS post_updated_at
FROM feeds f
JOIN posts p ON p.feed_id = f.id
ORDER BY f.name, p.published_at DESC;

-- name: DeleteStalePosts :exec
DELETE FROM posts
WHERE published_at < NOW() - INTERVAL '7 days';

-- name: GetPostsFromFollowedFeeds :many
SELECT 
    f.id AS feed_id,
    f.name AS feed_name,
    f.url AS feed_url,
    p.id AS post_id,
    p.title AS post_title,
    p.url AS post_url,
    p.description AS post_description,
    p.published_at AS post_published_at,
    p.created_at AS post_created_at,
    p.updated_at AS post_updated_at
FROM feeds f
JOIN posts p ON p.feed_id = f.id
JOIN feed_follows fo ON fo.feed_id = f.id
WHERE fo.user_id = $1
ORDER BY f.name, p.published_at DESC;
