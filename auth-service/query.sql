-- name: CreateUser :one
INSERT INTO users (id, name, email, password, policy)
VALUES ($1, $2, $3, $4, $5) RETURNING *;
SELECT * FROM users WHERE id = id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, email_verified = $4, password = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 RETURNING *;
SELECT * FROM users WHERE id = id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

