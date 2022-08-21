package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/uchennaemeruche/go-bank-api/api"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://postgres:postgres@localhost:5432/go_simple_bank?sslmode=disable"
	serverAddr = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to DB server: ", err)
	}
	store := db.NewStore(conn)

	server := api.NewServer(*store)

	server.Start(serverAddr)
}
