-- +migrate Up
create type  am_pm as ENUM ('am', 'pm');

CREATE TABLE turnip_prices
(
    id          BIGSERIAL PRIMARY KEY,
    discord_id  text      NOT NULL references accounts (discord_id),
    price       int       not null,
    am_pm    am_pm  not null,
    day_of_week int       not null,
    day_of_year int       not null,
    year        int       not null,
    created_at  timestamp not null default now()
);

create index on turnip_prices (discord_id);
create index on turnip_prices (created_at);
create index on turnip_prices (year, day_of_year, day_of_week);


create unique index
    on turnip_prices
        (
         discord_id,
         am_pm,
         day_of_week,
         day_of_year,
         year
            );

-- +migrate Down
drop table turnip_prices;
drop type am_pm;