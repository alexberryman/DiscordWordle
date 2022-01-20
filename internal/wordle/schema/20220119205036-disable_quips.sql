-- +migrate Up
create table disable_quips
(
    server_id text not null
);

create index on disable_quips(server_id);

-- +migrate Down
drop table disable_quips;
