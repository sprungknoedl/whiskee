-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table reviews (
	id bigserial primary key,
	user_id bigint not null references users(id),
	whisky_id bigint not null references whiskies(id),
	rating int not null,
	description text
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
drop table reviews
