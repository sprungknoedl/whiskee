-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table reviews add column date timestamp with time zone not null default now();


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table reviews drop column date;
