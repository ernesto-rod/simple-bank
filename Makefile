postgres:
	docker run --name=simple-bank-pg --publish=5433:5432 --env-file=./.dbenv --detach postgres:alpine

createdb:
	docker exec -it simple-bank-pg createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simple-bank-pg dropdb simple_bank

db_up:
	docker start simple-bank-pg

db_down:
	docker stop simple-bank-pg

db_connect:
		psql "postgresql://root:tr4nsactD3@localhost:5433/simple_bank"

migrate_up:
	migrate -path db/migration -database "postgresql://root:tr4nsactD3@localhost:5433/simple_bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://root:tr4nsactD3@localhost:5433/simple_bank?sslmode=disable" -verbose down

server:
	go run main.go

sqlc:
	sqlc generate

tests:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb db_up db_down migrate_up migrate_down server
