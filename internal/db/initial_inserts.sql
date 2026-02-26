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
    '$2a$12$Taa.42J00Xz8Jl3tpzNEEezhrFvZ1pRWUm4N4Tno43waWwo.DaG06',
    false,
    0,
    0,
    0,
    1,
    2
);
