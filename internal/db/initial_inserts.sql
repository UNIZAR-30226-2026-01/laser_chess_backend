INSERT INTO shop_item (
	price,
	level_requisite,
	item_type,
	is_default
)
VALUES (
    0,
    0,
    'board_skin',
    true
);

INSERT INTO shop_item (
	price,
	level_requisite,
	item_type,
	is_default
)
VALUES (
    0,
    0,
    'piece_skin',
    true
);

INSERT INTO public.account (
    mail,
    username,
    password_hash,
    is_deleted,
    "level",
    xp,
    "money",
    board_skin,
    piece_skin
)
VALUES(
    'user1@gmail.com',
    'user1', 
    '$2a$12$gZCIUMOOEuMhWvNlygdI3uTmXN90EEDBNVTCi/mWFuVfcRNQ.pIxi',
    false,
    0,
    0,
    0,
    1,
    2
), (
    'user2@gmail.com',
    'user2', 
    '$2a$12$Pm4HvVrVjvxfbihkZbfFneJtIZUb0G6CD7IERVrZQuaaAM/hASYT2',
    false,
    0,
    0,
    0,
    1,
    2
);

INSERT INTO item_owner (
    user_id,
    item_id
)
VALUES ( 1,1 ), ( 1,2 ), ( 2,1 ) , ( 2,2 );

INSERT INTO public."match" (
    p1_id,
    p2_id,
    p1_elo,
    p2_elo,
    "date",
    "winner",
    "termination",
    "match_type",
    board, movement_history,
    time_base, time_increment
) VALUES (
    1,
    2,
    1500,
    1600,
    '2026-02-22T15:04:05Z',
    'P1_WINS',
    'LASER',
    'RANKED',
    'ACE',
    'una ruta',
    300,
    5
);
