-- name: CreateLog :one
INSERT INTO log (id, message, thread_identifier, requestid, logtype, userid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;
SELECT * FROM log WHERE id = id;

-- name: GetLogByID :one
SELECT * FROM log WHERE id = $1;

-- name: GetLogByRequestID :many
SELECT * FROM log WHERE requestid = $1;

-- name: GetLogByThreadIdentifier :many
SELECT * FROM log WHERE thread_identifier = $1;

-- name: GetLogByLogType :many
SELECT * FROM log WHERE logtype = $1;

-- name: GetLogByUserID :many
SELECT * FROM log WHERE userid = $1;

-- name: DeleteLog :exec
DELETE FROM log WHERE id = $1;