// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createDevice = `-- name: CreateDevice :exec
insert into devices(user_id, device)
values($1, $2)
`

type CreateDeviceParams struct {
	UserID int64  `json:"user_id"`
	Device string `json:"device"`
}

func (q *Queries) CreateDevice(ctx context.Context, arg CreateDeviceParams) error {
	_, err := q.db.Exec(ctx, createDevice, arg.UserID, arg.Device)
	return err
}

const createUser = `-- name: CreateUser :exec
insert into users(created_at, email, phone_number, password_hash)
values(
        current_timestamp,
        $1,
        $2,
        $3
    )
`

type CreateUserParams struct {
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser, arg.Email, arg.PhoneNumber, arg.PasswordHash)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
delete from users
where id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getOneDevice = `-- name: GetOneDevice :one
select id, created_at, user_id, device
from devices
where devices.user_id = $1
`

func (q *Queries) GetOneDevice(ctx context.Context, userID int64) (*Devices, error) {
	row := q.db.QueryRow(ctx, getOneDevice, userID)
	var i Devices
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.Device,
	)
	return &i, err
}

const getUser = `-- name: GetUser :one
select id, created_at, updated_at, deleted_at, email, password_hash, phone_number, is_verified
from users
where id = $1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (*Users, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Email,
		&i.PasswordHash,
		&i.PhoneNumber,
		&i.IsVerified,
	)
	return &i, err
}

const listByVerification = `-- name: ListByVerification :many
select id, created_at, updated_at, deleted_at, email, password_hash, phone_number, is_verified
from users
where is_verified = $1
`

func (q *Queries) ListByVerification(ctx context.Context, isVerified pgtype.Bool) ([]*Users, error) {
	rows, err := q.db.Query(ctx, listByVerification, isVerified)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Users{}
	for rows.Next() {
		var i Users
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Email,
			&i.PasswordHash,
			&i.PhoneNumber,
			&i.IsVerified,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
select id, created_at, updated_at, deleted_at, email, password_hash, phone_number, is_verified
from users
`

func (q *Queries) ListUsers(ctx context.Context) ([]*Users, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Users{}
	for rows.Next() {
		var i Users
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Email,
			&i.PasswordHash,
			&i.PhoneNumber,
			&i.IsVerified,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateDevice = `-- name: UpdateDevice :exec
update devices
set device = $2
where user_id = $1
`

type UpdateDeviceParams struct {
	UserID int64  `json:"user_id"`
	Device string `json:"device"`
}

func (q *Queries) UpdateDevice(ctx context.Context, arg UpdateDeviceParams) error {
	_, err := q.db.Exec(ctx, updateDevice, arg.UserID, arg.Device)
	return err
}

const updateUser = `-- name: UpdateUser :exec
update users
set updated_at = current_timestamp,
    email = $2,
    phone_number = $3
where id = $1
`

type UpdateUserParams struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.ID, arg.Email, arg.PhoneNumber)
	return err
}

const userByEmail = `-- name: UserByEmail :one
SELECT is_verified,
    created_at,
    email,
    id,
    password_hash,
    phone_number
FROM users
WHERE email = $1
    AND "users"."deleted_at" IS NULL
ORDER BY "users"."id"
LIMIT 1
`

type UserByEmailRow struct {
	IsVerified   pgtype.Bool        `json:"is_verified"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	Email        string             `json:"email"`
	ID           int64              `json:"id"`
	PasswordHash string             `json:"password_hash"`
	PhoneNumber  string             `json:"phone_number"`
}

func (q *Queries) UserByEmail(ctx context.Context, email string) (*UserByEmailRow, error) {
	row := q.db.QueryRow(ctx, userByEmail, email)
	var i UserByEmailRow
	err := row.Scan(
		&i.IsVerified,
		&i.CreatedAt,
		&i.Email,
		&i.ID,
		&i.PasswordHash,
		&i.PhoneNumber,
	)
	return &i, err
}
