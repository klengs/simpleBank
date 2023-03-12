postgres:
	docker-compose up -d

createdb:
	docker-compose exec -it postgres createdb -U klengs simple_bank

dropdb: 
	docker-compose exec -it postgres dropdb -U klengs simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://klengs:Qwerty123@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://klengs:Qwerty123@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY:postgres createdb dropdb migrateup migratedown