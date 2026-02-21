erDiagram
item_owner ||--o{ account : is_owner
friendship ||--o{ account : is_friend
friendship ||--o{ account : is_friend
match ||--o{ account : player_one
match ||--o{ account : player_two
item_owner ||--o{ shop_item : is_owned
account }o--|| shop_item : equiped_board_skin
account }o--|| shop_item : equiped_piece_skin
ratings ||--|| account : elo

	item_owner {
		BIGINT user_id
		INTEGER item_id
	}

	account {
		BIGSERIAL account_id
		TEXT password_hash
		TEXT mail
		TEXT username
		BOOLEAN is_deleted
		INTEGER level
		INTEGER xp
		INTEGER money
		INTEGER board_skin
		INTEGER piece_skin
	}

	match {
		BIGSERIAL match_id
		BIGINT p1_id
		BIGINT p2_id
		INTEGER p1_elo
		INTEGER p2_elo
		TIMESTAMPTZ date
		WINNER winner
		TERMINATION termination
		MATCH_TYPE match_type
		BOARD_TYPE board
		TEXT movement_history
		INTEGER time_base
		INTEGER time_increment
	}

	friendship {
		BIGINT user1_id
		BIGINT user2_id
		BOOLEAN accepted_1
		BOOLEAN accepted_2
	}

	shop_item {
		SERIAL item_id
		INTEGER price
		INTEGER level_requisite
		ITEM_TYPE item_type
	}

	ratings {
		BIGINT user_id
		ELO_TYPE elo_type
		INTEGER value
	}
