begin;
insert into accounts (discord_id)
values ('1'),
       ('2'),
       ('3'),
       ('4'),
       ('5'),
       ('6');
insert into nicknames (discord_id, server_id, nickname)
VALUES ('1', '934007495737884712', '1'),
       ('2', '934007495737884712', '2'),
       ('3', '934007495737884712', '3'),
       ('4', '934007495737884712', '4'),
       ('5', '934007495737884712', '5'),
       ('6', '934007495737884712', '6');
insert into wordle_scores (discord_id, game_id, guesses)
VALUES ('1', 903, 1),
       ('1', 904, 1),
       ('1', 905, 1),
       ('1', 906, 1),
       ('1', 907, 1),
       ('1', 908, 1);
insert into wordle_scores (discord_id, game_id, guesses)
VALUES ('2', 903, 2),
       ('2', 904, 2),
       ('2', 905, 2),
       ('2', 906, 2),
       ('2', 907, 2);
insert into wordle_scores (discord_id, game_id, guesses)
VALUES ('3', 903, 3),
       ('3', 904, 3),
       ('3', 907, 3),
       ('3', 908, 3);
insert into wordle_scores (discord_id, game_id, guesses)
VALUES ('4', 903, 4),
       ('4', 907, 4),
       ('4', 908, 4);
insert into wordle_scores (discord_id, game_id, guesses)
VALUES ('5', 907, 5),
       ('5', 908, 5);
insert into wordle_scores (discord_id, game_id, guesses)
VALUES ('6', 908, 6);

commit;