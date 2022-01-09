-- +migrate Up
CREATE TABLE wordle_scores
(
    id          BIGSERIAL PRIMARY KEY,
    discord_id  text      NOT NULL references accounts (discord_id),
    game_id     int       not null,
    guesses        int       not null,
    created_at  timestamp not null default now()
);

create index on wordle_scores (discord_id);
create index on wordle_scores (created_at);
create index on wordle_scores (game_id);


create unique index
    on wordle_scores
    (
    discord_id,
    game_id
    );

-- +migrate Down
drop table wordle_scores;