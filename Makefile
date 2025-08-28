postgres:
	docker run --name tests-postgres-db -p 5432:5432 -e POSTGRES_USER=gustavo -e POSTGRES_PASSWORD=1910 -d postgres

createdb:
	docker exec -it tests-postgres-db --username=gustavo --owner=gustavo postgres

dropdb:
	docker exec -it tests-postgres-db dropdb postgres

migrateup:
	migrate -path db/migration -database "postgresql://gustavo:1910@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://gustavo:1910@localhost:5432/postgres?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./tests/... -cover

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test