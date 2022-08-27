postgres:
	docker run --name postgresdb -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres

createdb:
	docker exec -it postgresdb createdb --username=postgres --owner=postgres go_simple_bank

dropdb:
	docker exec -it postgresdb dropdb --username=postgres go_simple_bank

psql:
	docker exec -it postgresdb psql -U postgres

bash:
	docker exec -it postgresdb bash

migrateup:
	docker run -v /Users/emeruche/Practice/go-bank-api/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/go_simple_bank?sslmode=disable" up

migrateup_last:
	docker run -v /Users/emeruche/Practice/go-bank-api/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/go_simple_bank?sslmode=disable" up 1

migratedown:
	docker run -v /Users/emeruche/Practice/go-bank-api/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/go_simple_bank?sslmode=disable" down -all


migratedown_last:
	docker run -v /Users/emeruche/Practice/go-bank-api/db/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:postgres@localhost:5432/go_simple_bank?sslmode=disable" down 1

ci_migrateup:
	migrate  -path db/migrations/ -database "postgresql://postgres:postgres@localhost:5432/go_simple_bank?sslmode=disable" -verbose up

ci_migratedown:
	migrate  -path db/migrations/ -database "postgresql://postgres:postgres@localhost:5432/go_simple_bank?sslmode=disable" -verbose down -all


sqlc_init:
	/Users/emeruche/go/bin/sqlc init

sqlc_generate:
	/Users/emeruche/go/bin/sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go

mock:
	mockgen -destination db/mock/store.go --package mockdb github.com/uchennaemeruche/go-bank-api/db/sqlc Store

docker_build:
	docker build -t simplebank:latest .

docker_run:
	docker run --name simplebank -p 8080:8080 --network go-bank-api_app-tier -e GIN_MODE=release -e DB_SOURCE="postgresql://postgres:postgres@go-bank-api_postgresdb_1:5432/go_simple_bank?sslmode=disable" simplebank:latest

.PHONY:	postgres createdb dropdb psql bash migrateup migrateup_last migratedown migratedown_last sqlcinit sqlc_generate server mock docker_build docker_run