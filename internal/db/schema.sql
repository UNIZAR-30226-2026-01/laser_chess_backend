CREATE TYPE "winner" AS ENUM (
	'P1_WINS',
	'P2_WINS',
	'NONE'
);

CREATE TYPE "board_type" AS ENUM (
	'ACE',
	'CURIOSITY',
	'GRAIL',
	'MERCURY',
	'SOPHIE'
);

CREATE TYPE "match_type" AS ENUM (
	'PRIVATE',
	'RANKED',
	'BOTS'
);

CREATE TYPE "termination" AS ENUM (
	'OUT_OF_TIME',
	'SURRENDER',
	'LASER',
	'UNFINISHED',
	'DISCONNECTION'
);

CREATE TYPE "item_type" AS ENUM (
	'BOARD_SKIN',
	'PIECE_SKIN',
    'WIN_ANIMATION'
);

CREATE TYPE "elo_type" AS ENUM (
	'BLITZ',
	'EXTENDED',
	'RAPID',
	'CLASSIC'
);

CREATE TABLE IF NOT EXISTS "account" (
	"account_id" BIGSERIAL NOT NULL UNIQUE,
	"mail" VARCHAR(255) NOT NULL UNIQUE,
	"username" VARCHAR(50) NOT NULL UNIQUE,
	"password_hash" TEXT NOT NULL,

    -- params por defecto
	"is_deleted" BOOLEAN NOT NULL DEFAULT FALSE,
	"level" INTEGER NOT NULL DEFAULT 0,
	"xp" INTEGER NOT NULL DEFAULT 0,
	"money" INTEGER NOT NULL DEFAULT 0,

    -- items equipados
	"board_skin" INTEGER NOT NULL,
	"piece_skin" INTEGER NOT NULL,
	"win_animation" INTEGER NOT NULL,
	"avatar" INTEGER NOT NULL,

	PRIMARY KEY("account_id"),
	CHECK (
        "username" <> '' AND 
        "username" NOT LIKE '% %'
    ),
	CHECK ("mail" ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$')
);

CREATE TABLE IF NOT EXISTS "device" (
	"user_id" BIGSERIAL NOT NULL UNIQUE,
	"token" VARCHAR(255) NOT NULL UNIQUE,
	PRIMARY KEY("user_id", "token"),
	FOREIGN KEY("user_id") REFERENCES "account"("account_id")
);

CREATE TABLE IF NOT EXISTS "shop_item" (
	"item_id" SERIAL NOT NULL UNIQUE,
	"price" INTEGER NOT NULL,
	"level_requisite" INTEGER NOT NULL,
	"item_type" ITEM_TYPE NOT NULL,
	"is_default" BOOLEAN NOT NULL,
	PRIMARY KEY("item_id"),
	CHECK ("price" >= 0)
);

CREATE TABLE IF NOT EXISTS "item_owner" (
	"user_id" BIGINT NOT NULL,
	"item_id" INTEGER NOT NULL,
	PRIMARY KEY("user_id", "item_id"),
	FOREIGN KEY("user_id") REFERENCES "account"("account_id"),
	FOREIGN KEY("item_id") REFERENCES "shop_item"("item_id")
);

ALTER TABLE "account"
	ADD FOREIGN KEY("board_skin") REFERENCES "shop_item"("item_id");
ALTER TABLE "account"
	ADD FOREIGN KEY("piece_skin") REFERENCES "shop_item"("item_id");

CREATE TABLE IF NOT EXISTS "match" (
	"match_id" BIGSERIAL NOT NULL UNIQUE,
	"p1_id" BIGINT NOT NULL,
	"p2_id" BIGINT NOT NULL,
	"p1_elo" INTEGER NOT NULL,
	"p2_elo" INTEGER NOT NULL,
	"date" TIMESTAMPTZ NOT NULL,
	"winner" WINNER NOT NULL,
	"termination" TERMINATION NOT NULL,
	"match_type" MATCH_TYPE NOT NULL,
	"board" BOARD_TYPE NOT NULL,
	"movement_history" TEXT NOT NULL,
	"time_base" INTEGER NOT NULL,
	"time_increment" INTEGER NOT NULL,
	PRIMARY KEY("match_id"),
	FOREIGN KEY("p1_id") REFERENCES "account"("account_id"),
	FOREIGN KEY("p2_id") REFERENCES "account"("account_id"),
	CHECK ("p1_id" != "p2_id")
);

CREATE TABLE IF NOT EXISTS "friendship" (
	"user1_id" BIGINT NOT NULL,
	"user2_id" BIGINT NOT NULL,
	"accepted_1" BOOLEAN NOT NULL,
	"accepted_2" BOOLEAN NOT NULL,
	PRIMARY KEY("user1_id", "user2_id"),
	FOREIGN KEY("user1_id") REFERENCES "account"("account_id"),
	FOREIGN KEY("user2_id") REFERENCES "account"("account_id"),
	CHECK ("user1_id" < "user2_id")
);

CREATE TABLE IF NOT EXISTS "rating" (
	"user_id" BIGINT NOT NULL,
	"elo_type" ELO_TYPE NOT NULL,
	"value" INT NOT NULL DEFAULT 1500,
    "deviation" INT NOT NULL DEFAULT 350,
    "volatility" FLOAT NOT NULL DEFAULT 0.06,
    "last_updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),

	PRIMARY KEY("user_id", "elo_type"),
	FOREIGN KEY("user_id") REFERENCES "account"("account_id"),
	CHECK ("value" >= 0)
);

CREATE TABLE IF NOT EXISTS "refresh_session" (
    "session_id" BIGSERIAL PRIMARY KEY,
    "account_id" BIGINT NOT NULL,
    "token_hash" TEXT NOT NULL,
    "expires_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    FOREIGN KEY ("account_id") REFERENCES "account"("account_id")
);
-- btree para encontrar refresh token rapido
CREATE UNIQUE INDEX "refresh_session_token_hash_idx" ON "refresh_session" ("token_hash");
