
-- +migrate Up
alter table quips add column uses int not null default 0;

-- +migrate Down
alter table quips drop column uses;
