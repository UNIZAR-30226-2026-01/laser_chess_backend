INSERT INTO shop_item (
	price,
	level_requisite,
	item_type,
	is_default
)
VALUES (
    0,
    0,
    'BOARD_SKIN',
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
    'PIECE_SKIN',
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
    'WIN_ANIMATION',
    true
);

INSERT INTO public.account (
    account_id,
    mail,
    username,
    password_hash,
    is_deleted,
    level,
    xp,
    money,
    board_skin,
    piece_skin,
    win_animation,
    avatar
)
VALUES(
    1,
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
    2,
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
    3,
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
    4,
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
), (
    5,
    'user5@gmail.com',
    'user5', 
    '$2a$12$ucy8oE/C8vOeuBJP5tjRVOK803mdyPAChh1eRMV0tceWziPxiIewq',
    false,
    5,
    15,
    501,
    1,
    2,
    3,
    1
), (
    6,
    'user6@gmail.com',
    'user6', 
    '$2a$12$TJTTVIoZw09z9.Rre9ELdemP/a6JfAZ5DC3MvllTMbrMrJwflPKCK',
    false,
    6,
    16,
    601,
    1,
    2,
    3,
    1
), (
    7,
    'user7@gmail.com',
    'user7', 
    '$2a$12$iiEC/A6ZbKgeh75egfPlKupOQW7JL5VPXz.IAlvJ6lO3oUUAtjFWe',
    false,
    7,
    17,
    701,
    1,
    2,
    3,
    1
), (
    8,
    'user8@gmail.com',
    'user8', 
    '$2a$12$u548i67ylzsb1vOykpWUouOxH6wdEaO1J.tCXtaFGyNDF4rOIDV7W',
    false,
    8,
    18,
    801,
    1,
    2,
    3,
    1
), (
    9,
    'user9@gmail.com',
    'user9', 
    '$2a$12$R77H7YqkOPEGQep4n7oHUeFx2sqG0msqYT3uN.2yWPdM8copfY5EK',
    false,
    9,
    19,
    901,
    1,
    2,
    3,
    1
), (
    10,
    'user10@gmail.com',
    'user10', 
    '$2a$12$dVqxoDl/Aubzgd/uIpQLvuAAqyXg3KvhpFnNj1GHDEV.Whvb/gABW',
    false,
    10,
    20,
    1001,
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
    'BLITZ',
    1500
), (
    1,
    'RAPID',
    4500
), (
    1,
    'CLASSIC',
    1500
), (
    1,
    'EXTENDED',
    1500
), (
    2,
    'BLITZ',
    1600
), (
    2,
    'RAPID',
    1600
), (
    2,
    'CLASSIC',
    4600
), (
    2,
    'EXTENDED',
    1600
), (
    3,
    'BLITZ',
    4700
), (
    3,
    'RAPID',
    1700
), (
    3,
    'CLASSIC',
    1700
), (
    3,
    'EXTENDED',
    1700
), (
    4,
    'BLITZ',
    1800
), (
    4,
    'RAPID',
    1800
), (
    4,
    'CLASSIC',
    1800
),(
    4,
    'EXTENDED',
    4800
), (
    5,
    'BLITZ',
    4900
), (
    5,
    'RAPID',
    1900
), (
    5,
    'CLASSIC',
    1900
), (
    5,
    'EXTENDED',
    1900
), (
    6,
    'BLITZ',
    2000
), (
    6,
    'RAPID',
    4000
), (
    6,
    'CLASSIC',
    2000
), (
    6,
    'EXTENDED',
    2000
), (
    7,
    'BLITZ',
    2100
), (
    7,
    'RAPID',
    2100
), (
    7,
    'CLASSIC',
    4100
), (
    7,
    'EXTENDED',
    2100
), (
    8,
    'BLITZ',
    2200
), (
    8,
    'RAPID',
    2200
), (
    8,
    'CLASSIC',
    2200
), (
    8,
    'EXTENDED',
    4200
),(
    9,
    'BLITZ',
    2300
), (
    9,
    'RAPID',
    4300
), (
    9,
    'CLASSIC',
    2300
), (
    9,
    'EXTENDED',
    2300
), (
    10,
    'BLITZ',
    2400
), (
    10,
    'RAPID',
    4400
), (
    10,
    'CLASSIC',
    2400
), (
    10,
    'EXTENDED',
    2400
);

INSERT INTO shop_item (
    price, level_requisite, item_type, is_default
) VALUES (
    50, 0, 'BOARD_SKIN', false
), (
    1000, 0, 'BOARD_SKIN', false
), (
    1, 10, 'BOARD_SKIN', false
);

SELECT setval(
    pg_get_serial_sequence('public.account', 'account_id'),
    COALESCE((SELECT MAX(account_id) FROM public.account), 1)
);