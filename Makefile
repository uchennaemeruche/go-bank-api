postgres:
	docker run --name postgresdb -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres

createdb:
	docker exec -it postgresdb createdb --username=postgres --owner=postgres go_simple_bank

dropdb:
	docker exec -it postgresdb dropdb --username=postgres go_simple_bank
