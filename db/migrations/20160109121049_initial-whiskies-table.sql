-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table whiskies (
	id bigserial primary key,
	type text not null,
	distillery text not null,
	name text,
	age integer not null,
	abv double precision,
	description text,

	picture text,
	thumbnail text,

	rating double precision not null default 0.0,
	ratings integer not null default 0,
	reviews integer not null default 0
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table whiskies
