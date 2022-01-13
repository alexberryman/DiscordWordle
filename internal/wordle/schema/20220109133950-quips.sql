
-- +migrate Up
create table quips (
    id BIGSERIAL PRIMARY KEY,
    score_value int not null,
    quip text not null,
    inside_joke bool not null default false,
    inside_joke_server_id text,
    created_by_account text not null references accounts,
    created_at timestamp not null default now(),
    CONSTRAINT if_inside_joke_then_inside_joke_server_id_is_not_null
        CHECK ( (NOT inside_joke) OR (quips.inside_joke_server_id IS NOT NULL) ),
    CONSTRAINT no_empty_quips
        CHECK ( quip != '' ),
    CONSTRAINT no_empty_server_id
        CHECK ( quips.inside_joke_server_id != '' )
);

INSERT INTO public.accounts (discord_id, time_zone) VALUES ('229048840454406145', 'America/Chicago');
insert into quips (score_value, quip, inside_joke, created_by_account)
    VALUES (0, 'woof 😨', false, '229048840454406145');
insert into quips (score_value, quip, inside_joke, created_by_account)
    VALUES (1, 'lol, try hard', false, '229048840454406145');
insert into quips (score_value, quip, inside_joke, created_by_account)
    VALUES (2, '2? Dope. Almost like two pairs of Air Force Ones', false, '229048840454406145');
insert into quips (score_value, quip, inside_joke, created_by_account)
    VALUES (3, 'Three is perfectly reasonable', false, '229048840454406145');
insert into quips (score_value, quip, inside_joke, inside_joke_server_id,created_by_account)
    VALUES (3, 'Nice, that''s what Brad got that one time which was par for the course', true, '700123904815136879', '229048840454406145');
insert into quips (score_value, quip, inside_joke, created_by_account)
    VALUES (4, '4 is the new 6. You should have just given up after 3.', false, '229048840454406145');
insert into quips (score_value, quip, inside_joke, created_by_account)
    VALUES (5, 'Remember when Arby''s did the 5 for 5 deal? Good times.', false, '229048840454406145');
insert into quips (score_value, quip, inside_joke, created_by_account)
    VALUES (6, '6? I bet you waited to the last morning to do your homework too.', false, '229048840454406145');

create index on quips (score_value);
create index on quips (score_value, inside_joke, inside_joke_server_id);

-- +migrate Down
DROP TABLE quips;