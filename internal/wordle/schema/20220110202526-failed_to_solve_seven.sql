
-- +migrate Up
update wordle_scores set guesses = 7 where guesses = 0;
-- +migrate Down
update wordle_scores set guesses = 0 where guesses = 7;
