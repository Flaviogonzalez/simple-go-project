-- name: CreateSession :one
INSERT INTO sessions (id, session_token, csrf_token, user_id, expires) VALUES ($1, $2, $3, $4, $5) RETURNING *;
SELECT * FROM sessions WHERE id = id;

-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;