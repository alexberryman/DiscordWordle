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