-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table users (
	id bigserial primary key,
	auth0 text unique,
	name text,
	email text,
	picture text
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table users;
