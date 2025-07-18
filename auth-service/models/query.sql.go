// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package models

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, name, email, password, policy)
VALUES ($1, $2, $3, $4, $5) RETURNING id, name, email, email_verified, policy, password, created_at, updated_at
`

type CreateUserParams struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
	Policy   int32
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Policy,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.EmailVerified,
		&i.Policy,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, email_verified, policy, password, created_at, updated_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.EmailVerified,
		&i.Policy,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, email, email_verified, policy, password, created_at, updated_at FROM users WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.EmailVerified,
		&i.Policy,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, email, email_verified, policy, password, created_at, updated_at FROM users
ORDER BY created_at DESC
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.EmailVerified,
			&i.Policy,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, email_verified = $4, password = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 RETURNING id, name, email, email_verified, policy, password, created_at, updated_at
`

type UpdateUserParams struct {
	ID            uuid.UUID
	Name          string
	Email         string
	EmailVerified sql.NullTime
	Password      string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.EmailVerified,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.EmailVerified,
		&i.Policy,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
