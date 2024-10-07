-- name: CreateRemainder :one
insert into remainder(id, created_at, updated_at, subject, description, has_priority, timing, userid)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetRemaindersByUser :many
select * from remainder where userid = $1;

-- name: GetRemainderByID :one
select * from remainder where id = $1;

-- name: UpdateRemainder :one
update remainder
set subject = $1,
description = $2,
has_priority = $3,
timing = $4
where id = $5 and userid = $6
returning *;

-- name: DeleteRemainder :one
delete from remainder where id = $1 and userid = $2
returning *;