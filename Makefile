include .env
export

postgres:
	docker run --name tests-postgres-db -p 5432:5432 \
		-e POSTGRES_USER=gustavo \
		-e POSTGRES_PASSWORD=1910 \
		-d postgres:15-alpine

createdb:
	docker exec -it tests-postgres-db createdb --username=gustavo --owner=gustavo postgres

dropdb:
	docker exec -it tests-postgres-db dropdb postgres

migrateup:
	migrate -path db/migration -database "$(DB_DSN)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_DSN)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./tests/... -cover

server:
	DB_DSN=$(DB_DSN) go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
