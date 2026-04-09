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
    win_animation,
    avatar
)
VALUES(
    'user1@gmail.com',
    'user1', 
    '$2a$12$Q0RsprBAXFgSSttIet3nNe/rwKOFOJychVd9F2BH6q/Q3Pp5lnV3.',
    false,
    1,
    11,
    101,
    1,
    2,
    3,
    1
), (
    'user2@gmail.com',
    'user2', 
    '$2a$12$ApqnzgaPJ3LpimwZGjxJk.rBhWQ5EIQ.YNBDxXJZKmDXlxY42cNFK',
    false,
    2,
    12,
    201,
    1,
    2,
    3,
    1
), (
    'user3@gmail.com',
    'user3', 
    '$2a$12$HX2kns7L6joaJo07PrGafO8Sjz044snkwRBIh7pjdtHps4u.2kBLa',
    false,
    3,
    13,
    301,
    1,
    2,
    3,
    1
), (
    'user4@gmail.com',
    'user4', 
    '$2b$12$y4/FIcy88bFV6JArmNtq6.ovknJ3I7ynqepkQ6s/4usAtcjH3ouwC',
    false,
    4,
    14,
    401,
    1,
    2,
    3,
    1
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
    '',
    300,
    5
), (
    1,
    2,
    1500,
    1600,
    '2026-02-22T15:04:05Z',
    'NONE',
    'UNFINISHED',
    'RANKED',
    'CURIOSITY',
    'Rf1%j1,j4,i4,i5,j5,j9%{300};Tg6:f6%a8,a5,b5,b4,a4,a0%{300};Rb4%j1,j4,i4,i5,j5,j9%{295};Ri5xf6%a8,a5,b5,b4,e4,e5,f5,f6%{290};',
    300,
    0
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
), (
    1,
    3,
    false,
    true
), (
    1,
    4,
    true,
    false
);

INSERT INTO public."rating" (
    user_id,
    elo_type,
    value
) VALUES (
    1,
    'blitz',
    1
), (
    1,
    'rapid',
    11
), (
    1,
    'classic',
    111
), (
    1,
    'extended',
    1111
), (
    2,
    'blitz',
    2
), (
    2,
    'rapid',
    22
), (
    2,
    'classic',
    222
), (
    2,
    'extended',
    2222
);

INSERT INTO shop_item (
    price, level_requisite, item_type, is_default
) VALUES (
    50, 0, 'board_skin', false
), (
    1000, 0, 'board_skin', false
), (
    1, 10, 'board_skin', false
);