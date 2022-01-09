-- +migrate Up
-- +migrate StatementBegin
DO $$
    BEGIN
        IF EXISTS
            ( SELECT 1
              FROM   information_schema.tables
              WHERE  table_schema = 'public'
                AND    table_name = 'users'
            )
        THEN
            insert into accounts (discord_id, time_zone)
            select discord_id, time_zone
            from users;
        END IF ;
    END
$$ ;
-- +migrate StatementEnd

-- +migrate StatementBegin
DO $$
    BEGIN
        IF EXISTS
            ( SELECT 1
              FROM   information_schema.tables
              WHERE  table_schema = 'public'
                AND    table_name = 'server_context'
            )
        THEN
            insert into nicknames (discord_id, server_id, nickname)
            select discord_id, server_id, username
            from server_context;
        END IF ;
    END
$$ ;
-- +migrate StatementEnd

-- +migrate StatementBegin
DO $$
    BEGIN
        IF EXISTS
            ( SELECT 1
              FROM   information_schema.tables
              WHERE  table_schema = 'public'
                AND    table_name = 'prices'
            )
        THEN
            insert into turnip_prices (discord_id, price, am_pm, day_of_week, day_of_year, year)
            SELECT discord_id, price, meridiem::text::am_pm, day_of_week, day_of_year, year
            from prices;
        END IF ;
    END
$$ ;
-- +migrate StatementEnd


-- +migrate Down
truncate turnip_prices, nicknames, accounts CASCADE;