include .env
export

BINARY_NAME=main
MIGRATION_PATH=internal/db/migration

postgres:
	docker compose up -d

migrate-up: postgres
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose up

sqlc: migrate-up
	sqlc generate

build: sqlc
	go build -o bin/$(BINARY_NAME) cmd/main.go

run: build
	./bin/$(BINARY_NAME)

clean:
	rm -rf bin/


# Para hacer cualquier modificacion a las tablas
# Uso: make new-migration name=<el nombre que sea>
new-migration:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)





