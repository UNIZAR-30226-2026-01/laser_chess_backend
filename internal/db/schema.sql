CREATE TYPE "winner" AS ENUM (
	'P1_WINS',
	'P2_WINS',
	'DRAW',
	'NONE'
);

CREATE TYPE "board_type" AS ENUM (
	'ACE',
	'CURIOSITY',
	'SOPHIE',
	'GRAIL',
	'MERCURY'
);

CREATE TYPE "match_type" AS ENUM (
	'RANKED',
	'FRIENDLY',
	'PRIVATE',
	'BOTS'
);

CREATE TYPE "termination" AS ENUM (
	'OUT_OF_TIME',
	'SURRENDER',
	'LASER',
	'UNFINISHED'
);

CREATE TYPE "item_type" AS ENUM (
	'board_skin',
	'piece_skin'
);

CREATE TYPE "elo_type" AS ENUM (
	'blitz',
	'bullet',
	'rapid',
	'classic'
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
	PRIMARY KEY("account_id")
);

CREATE TABLE IF NOT EXISTS "shop_item" (
	"item_id" SERIAL NOT NULL UNIQUE,
	"price" INTEGER NOT NULL,
	"level_requisite" INTEGER NOT NULL,
	"item_type" ITEM_TYPE NOT NULL,
	PRIMARY KEY("item_id")
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
	FOREIGN KEY("p2_id") REFERENCES "account"("account_id")
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
	"value" INT NOT NULL,
	PRIMARY KEY("user_id", "elo_type"),
	FOREIGN KEY("user_id") REFERENCES "account"("account_id")
);



