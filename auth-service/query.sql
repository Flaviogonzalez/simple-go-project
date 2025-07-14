-- name: CreateUser :execresult
INSERT INTO users (id, name, email, password, policy)
VALUES (?, ?, ?, ?, ?);
SELECT * FROM users WHERE id = id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY createdAt DESC;

-- name: UpdateUser :execresult
UPDATE users
SET name = name, email = email, emailVerified = emailVerified, password = password, updateAt = CURRENT_TIMESTAMP
WHERE id = id;
SELECT * FROM users WHERE id = id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = id;

