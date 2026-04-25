-include .env
export

BINARY_NAME=main
DB_CONTAINER=postgres_db
SCHEMA_FILE=internal/db/schema.sql
INSERTS_FILE=internal/db/initial_inserts.sql

postgres:
	docker compose up -d

apply-schema: postgres
	docker exec -i $(DB_CONTAINER) psql "$(DB_URL)" < $(SCHEMA_FILE)

seed:
	go run cmd/seed/main.go

seed-debug:
	SEED_DEBUG=true go run cmd/seed/main.go

sqlc: 
	sqlc generate

build: sqlc
	go build -o bin/$(BINARY_NAME) cmd/main.go

run: 
	./bin/$(BINARY_NAME)

clean:
	rm -rf bin/

reset-db: postgres
	docker exec -i $(DB_CONTAINER) psql "$(DB_URL)" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	make apply-schema
	make seed

reset-db-debug: postgres
	docker exec -i $(DB_CONTAINER) psql "$(DB_URL)" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	make apply-schema
	make seed-debug
