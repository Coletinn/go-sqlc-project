-- Users queries
-- name: GetUserByID :one
SELECT id, name, email, phone, created_at FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, phone) 
VALUES ($1, $2, $3) 
RETURNING id, name, email, phone, created_at;

-- name: ListUsers :many
SELECT id, name, email, phone, created_at FROM users 
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users 
SET name = $2, email = $3, phone = $4 
WHERE id = $1 
RETURNING id, name, email, phone, created_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;