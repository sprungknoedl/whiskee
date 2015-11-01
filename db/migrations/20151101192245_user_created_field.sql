
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table users
  add column created timestamp with time zone default now();


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table users drop column created;
