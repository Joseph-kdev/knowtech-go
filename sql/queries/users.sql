-- name: AddUserToDatabase :one
INSERT INTO users(id, email, created_at, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (email) DO NOTHING
RETURNING *;