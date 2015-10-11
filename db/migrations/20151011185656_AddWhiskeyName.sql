-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table whiskeys add column name text not null default '';


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table whiskeys drop column name;
