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

migrate:
	docker exec -it postgresdb psql -U postgres --dbname=go_simple_bank --file=/Users/emeruche/Practice/go-bank-api/db/migrations/000001_init_schema.up.sql 

