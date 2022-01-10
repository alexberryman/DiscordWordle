
-- +migrate Up
create table responses (
    id BIGSERIAL PRIMARY KEY,
    score_value int not null,
    inject_score bool not null default false,
    response text not null,
    inside_joke bool not null default false,
    inside_joke_server_id text,
    created_by_account text not null references accounts,
    created_at timestamp not null default now(),
    CONSTRAINT if_inside_joke_then_inside_joke_server_id_is_not_null
        CHECK ( (NOT inside_joke) OR (responses.inside_joke_server_id IS NOT NULL) ),
    CONSTRAINT no_empty_responses
        CHECK ( response != '' )
);

insert into responses (score_value, inject_score, response, inside_joke, created_by_account)
    VALUES (0, false, 'woof 😨', false, '229048840454406145');
insert into responses (score_value, inject_score, response, inside_joke, created_by_account)
    VALUES (1, false, 'lol, try hard', false, '229048840454406145');
insert into responses (score_value, inject_score, response, inside_joke, created_by_account)
    VALUES (2, true, '%d? Dope. Almost like two pairs of Air Force Ones', false, '229048840454406145');
insert into responses (score_value, inject_score, response, inside_joke, created_by_account)
    VALUES (3, false, 'Three is perfectly reasonable', false, '229048840454406145');
insert into responses (score_value, inject_score, response, inside_joke, inside_joke_server_id,created_by_account)
    VALUES (3, false, 'Nice, that''s what Brad got that one time which was par for the course', true, '700123904815136879', '229048840454406145');
insert into responses (score_value, inject_score, response, inside_joke, created_by_account)
    VALUES (4, true, '%d is the new 6. You should have just given up after 3.', false, '229048840454406145');
insert into responses (score_value, inject_score, response, inside_joke, created_by_account)
    VALUES (5, false, 'Remember when Arby''s did the 5 for 5 deal? Good times.', false, '229048840454406145');
insert into responses (score_value, inject_score, response, inside_joke, created_by_account)
    VALUES (6, false, '6? I bet you waited to the last morning to do your homework too.', false, '229048840454406145');

create index on responses (score_value);
create index on responses (score_value, inside_joke, inside_joke_server_id);

-- +migrate Down
DROP TABLE responses;