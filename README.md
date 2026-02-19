# Laser Chess - Backend [![CI](https://github.com/UNIZAR-30226-2026-01/laser_chess_backend/actions/workflows/ci.yaml/badge.svg)](https://github.com/UNIZAR-30226-2026-01/laser_chess_backend/actions/workflows/ci.yaml)

Backend del juego de mesa online Laser Chess.

## Diseño preliminar de la base de datos

```mermaid
erDiagram
	user ||--o{ match : player1
	user ||--o{ match : player2
	user ||--o{ friendship : sender
	user ||--o{ friendship : reciever
	user ||--o{ item_owner : is_owner
	shop_item ||--o{ item_owner : is_owned

	item_owner {
		INTEGER owner_id
		INTEGER owned_id
	}

	user {
		INTEGER user_id
		TEXT password_hash
		TEXT mail
		TEXT username
		TEXT acount_state
		INTEGER elo
		INTEGER level
		INTEGER xp
		INTEGER money
	}

	match {
		INTEGER match_id
		INTEGER p1_id
		INTEGER p2_id
		INTEGER p1_elo
		INTEGER p2_elo
		TIMESTAMP date
		WINNER winner
		TERMINATION termination
		MATCH_TYPE public_ranked
		BOARD_TYPE board
		TEXT movement_history
		INTEGER time_base
		INTEGER time_increment
	}

	friendship {
		INTEGER user1_id
		INTEGER user2_id
		BOOLEAN accepted
	}

	shop_item {
		INTEGER item_id
		INTEGER price
		INTEGER level_requisite
	}
```
---
```mermaid

classDiagram
    class board_type {
        <<enumeration>>
        ACE
        CURIOSITY
		GRAIL
		SOPHIE
        MERCURY
    }

	class winner {
        <<enumeration>>
		P1_WINS
		P2_WINS
		DRAW
		NONE
    }

	class match_type {
        <<enumeration>>
		RANKED
		FRIENDLY
		PRIVATE
		BOTS
    }

	class termination {
        <<enumeration>>
		TIME
		SURRENDER
		LASER
		UNFINISHED
    }
	