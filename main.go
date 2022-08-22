package main

import (
	"database/sql"
	"log"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"github.com/uchennaemeruche/go-bank-api/api"
	db "github.com/uchennaemeruche/go-bank-api/db/sqlc"
	"github.com/uchennaemeruche/go-bank-api/util"
)

func main() {

	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal("Could not load configuration file: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB server: ", err)
	}
	store := db.NewStore(conn)

	server := api.NewServer(store)

	server.Start(config.ServerAddress)
}
