-- +migrate Up
alter table turnip_prices
    add column week int;

update turnip_prices
set week = newdata.w
from (select extract(WEEK from created_at at time zone 'utc' at time zone a.time_zone) w, a.discord_id id
      from turnip_prices
               inner join accounts a on turnip_prices.discord_id = a.discord_id
     ) newdata
where turnip_prices.week is null
  and turnip_prices.discord_id = newdata.id;

alter table turnip_prices
    alter column week set not null;


create unique index id_date_price_idx
    on turnip_prices
        (
         discord_id,
         am_pm,
         day_of_week,
         day_of_year,
         year,
         week
            );

-- +migrate Down
drop index id_date_price_idx;

alter table turnip_prices
    drop column week;

