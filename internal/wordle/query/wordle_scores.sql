-- name: GetScoreHistoryByAccount :many
SELECT *
FROM wordle_scores
         inner join nicknames nick on wordle_scores.discord_id = nick.discord_id
WHERE nick.discord_id = $1
  and nick.server_id = $2
order by game_id;

-- name: ListScores :many
SELECT *
FROM wordle_scores
ORDER BY created_at;

-- name: CountScoresByDiscordId :one
SELECT count(*)
FROM wordle_scores
where discord_id = $1;

-- name: CreateScore :one
INSERT INTO wordle_scores (discord_id, game_id, guesses)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateScore :one
update wordle_scores
set guesses = $2
where discord_id = $1
  and game_id = $3
returning *;

-- name: DeleteScoresForUser :exec
DELETE
FROM wordle_scores
WHERE discord_id = $1;

-- name: GetScoresByServerId :many
with max_game_week as (select max(game_id / 7) game_week
                       from wordle_scores
                                inner join nicknames n2 on wordle_scores.discord_id = n2.discord_id
                       where n2.server_id = $1
)
select n.nickname,
       json_agg(guesses order by s.game_id)             guesses_per_game,
       json_agg((7 - s.guesses) ^ 2 order by s.game_id) points_per_game,
       sum((7 - s.guesses) ^ 2)                         total
from wordle_scores s
         inner join nicknames n on s.discord_id = n.discord_id
         inner join max_game_week g on g.game_week = s.game_id / 7
where n.server_id = $1
group by n.nickname
order by sum((7 - s.guesses) ^ 2) desc;