-- +migrate Up
create table nicknames
(
    discord_id text not null references accounts (discord_id),
    server_id  text not null,
    nickname   text not null,
    primary key (discord_id,server_id)
);

create index on nicknames (server_id);
create index on nicknames (discord_id);
create unique index on nicknames (server_id, discord_id);

-- +migrate Down
drop table nicknames;