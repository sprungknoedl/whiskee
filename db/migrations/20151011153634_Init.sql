-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table users (
	id text primary key,
	email text not null,
	gravatar text not null
);

create table whiskeys (
	id bigserial primary key,
	distillery text not null,
	type text not null,
	age integer not null,
	abv double precision not null,
	size double precision not null
);

create table posts (
	id bigserial primary key,
	date timestamp with time zone not null,
	body text not null,
	user_id text not null,
	whiskey_id bigint not null,
	security text not null default 'public',
	foreign key (user_id) references users (id),
	foreign key (whiskey_id) references whiskeys (id)
);

create table friends (
	a text,
	b text,
	primary key (a, b),
	foreign key (a) references users (id),
	foreign key (b) references users (id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table friends;
drop table posts;
drop table whiskeys;
drop table users;
