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

apply-inserts: postgres
	docker exec -i $(DB_CONTAINER) psql "$(DB_URL)" < $(INSERTS_FILE)

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
	make apply-inserts
