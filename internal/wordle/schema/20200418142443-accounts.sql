-- +migrate Up
CREATE TABLE accounts
(
    discord_id  text unique not null PRIMARY KEY,
    time_zone   text        not null default 'America/Chicago'
);

create index on accounts (discord_id);;

-- +migrate Down
DROP TABLE accounts;