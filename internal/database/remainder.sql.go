// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: remainder.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRemainder = `-- name: CreateRemainder :one
insert into remainder(id, created_at, updated_at, subject, description, has_priority, timing, userid)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning id, created_at, updated_at, subject, description, has_priority, timing, userid
`

type CreateRemainderParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Subject     string
	Description string
	HasPriority bool
	Timing      time.Time
	Userid      uuid.UUID
}

func (q *Queries) CreateRemainder(ctx context.Context, arg CreateRemainderParams) (Remainder, error) {
	row := q.db.QueryRowContext(ctx, createRemainder,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Subject,
		arg.Description,
		arg.HasPriority,
		arg.Timing,
		arg.Userid,
	)
	var i Remainder
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Subject,
		&i.Description,
		&i.HasPriority,
		&i.Timing,
		&i.Userid,
	)
	return i, err
}

const deleteRemainder = `-- name: DeleteRemainder :one
delete from remainder where id = $1 and userid = $2
returning id, created_at, updated_at, subject, description, has_priority, timing, userid
`

type DeleteRemainderParams struct {
	ID     uuid.UUID
	Userid uuid.UUID
}

func (q *Queries) DeleteRemainder(ctx context.Context, arg DeleteRemainderParams) (Remainder, error) {
	row := q.db.QueryRowContext(ctx, deleteRemainder, arg.ID, arg.Userid)
	var i Remainder
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Subject,
		&i.Description,
		&i.HasPriority,
		&i.Timing,
		&i.Userid,
	)
	return i, err
}

const getRemainderByID = `-- name: GetRemainderByID :one
select id, created_at, updated_at, subject, description, has_priority, timing, userid from remainder where id = $1
`

func (q *Queries) GetRemainderByID(ctx context.Context, id uuid.UUID) (Remainder, error) {
	row := q.db.QueryRowContext(ctx, getRemainderByID, id)
	var i Remainder
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Subject,
		&i.Description,
		&i.HasPriority,
		&i.Timing,
		&i.Userid,
	)
	return i, err
}

const getRemaindersByUser = `-- name: GetRemaindersByUser :many
select id, created_at, updated_at, subject, description, has_priority, timing, userid from remainder where userid = $1
`

func (q *Queries) GetRemaindersByUser(ctx context.Context, userid uuid.UUID) ([]Remainder, error) {
	rows, err := q.db.QueryContext(ctx, getRemaindersByUser, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Remainder
	for rows.Next() {
		var i Remainder
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Subject,
			&i.Description,
			&i.HasPriority,
			&i.Timing,
			&i.Userid,
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

const updateRemainder = `-- name: UpdateRemainder :one
update remainder
set subject = $1,
description = $2,
has_priority = $3,
timing = $4
where id = $5 and userid = $6
returning id, created_at, updated_at, subject, description, has_priority, timing, userid
`

type UpdateRemainderParams struct {
	Subject     string
	Description string
	HasPriority bool
	Timing      time.Time
	ID          uuid.UUID
	Userid      uuid.UUID
}

func (q *Queries) UpdateRemainder(ctx context.Context, arg UpdateRemainderParams) (Remainder, error) {
	row := q.db.QueryRowContext(ctx, updateRemainder,
		arg.Subject,
		arg.Description,
		arg.HasPriority,
		arg.Timing,
		arg.ID,
		arg.Userid,
	)
	var i Remainder
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Subject,
		&i.Description,
		&i.HasPriority,
		&i.Timing,
		&i.Userid,
	)
	return i, err
}
