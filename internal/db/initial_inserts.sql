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

INSERT INTO shop_item (
	price,
	level_requisite,
	item_type,
	is_default
)
VALUES (
    0,
    0,
    'win_animation',
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
    piece_skin,
    win_animation
)
VALUES(
    'user1@gmail.com',
    'user1', 
    '$2a$12$Q0RsprBAXFgSSttIet3nNe/rwKOFOJychVd9F2BH6q/Q3Pp5lnV3.',
    false,
    0,
    0,
    0,
    1,
    2,
    3
), (
    'user2@gmail.com',
    'user2', 
    '$2a$12$ApqnzgaPJ3LpimwZGjxJk.rBhWQ5EIQ.YNBDxXJZKmDXlxY42cNFK',
    false,
    0,
    0,
    0,
    1,
    2,
    3
), (
    'user3@gmail.com',
    'user3', 
    '$2a$12$HX2kns7L6joaJo07PrGafO8Sjz044snkwRBIh7pjdtHps4u.2kBLa',
    false,
    0,
    0,
    0,
    1,
    2,
    3
);

INSERT INTO item_owner (
    user_id,
    item_id
)
VALUES ( 1,1 ), ( 1,2 ), ( 1,3 ), ( 2,1 ) , ( 2,2 ), ( 2,3 );

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

INSERT INTO public."friendship" (
    user1_id,
    user2_id,
    accepted_1,
    accepted_2
) VALUES (
    1,
    2,
    true,
    true
);