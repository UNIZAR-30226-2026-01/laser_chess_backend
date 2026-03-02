```mermaid
erDiagram
    account ||--o{ item_owner : owns
    account ||--o{ match : plays_as_p1
    account ||--o{ match : plays_as_p2
    account ||--o{ friendship : friend_1
    account ||--o{ friendship : friend_2
    account ||--o{ rating : has_rating
    account ||--o{ refresh_session : has_session
    account }o--|| shop_item : equips_board
    account }o--|| shop_item : equips_piece
    shop_item ||--o{ item_owner : is_owned

    account {
        BIGSERIAL account_id PK
        VARCHAR mail UK
        VARCHAR username UK
        TEXT password_hash
        BOOLEAN is_deleted
        INTEGER level
        INTEGER xp
        INTEGER money
        INTEGER board_skin FK
        INTEGER piece_skin FK
    }

    shop_item {
        SERIAL item_id PK
        INTEGER price
        INTEGER level_requisite
        ITEM_TYPE item_type
        BOOLEAN is_default
    }

    item_owner {
        BIGINT user_id PK, FK
        INTEGER item_id PK, FK
    }

    match {
        BIGSERIAL match_id PK
        BIGINT p1_id FK
        BIGINT p2_id FK
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
        BIGINT user1_id PK, FK
        BIGINT user2_id PK, FK
        BOOLEAN accepted_1
        BOOLEAN accepted_2
    }

    rating {
        BIGINT user_id PK, FK
        ELO_TYPE elo_type PK
        INTEGER value
    }

    refresh_session {
        BIGSERIAL session_id PK
        BIGINT account_id FK
        TEXT token_hash UK
        TIMESTAMPTZ expires_at
        TIMESTAMPTZ created_at
    }

