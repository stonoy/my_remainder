-- +goose Up
create table remainder(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    subject text not null,
    description text not null,
    has_priority boolean not null,
    timing timestamp not null,
    userid uuid not null
    references users(id)
    on delete cascade
);

-- +goose Down
drop table remainder;