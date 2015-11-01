-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table users 
	drop column gravatar,
	add column name text not null,
	add column nick text not null,
	add column picture text not null;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table users
	drop column name,
	drop column nick,
	drop column picture,
	add column gravatar text not null;
