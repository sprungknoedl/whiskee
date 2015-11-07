
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
alter table whiskeys
  add column picture text not null default '/static/img/default.jpg',
  add column thumb text not null default '/static/img/default-thumb.jpg';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
alter table whiskeys
  drop column picture,
  drop column thumb;
