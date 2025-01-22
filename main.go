package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/shreyanshsharma88/golang-bank/api"
	db "github.com/shreyanshsharma88/golang-bank/db/sqlc"
	"github.com/shreyanshsharma88/golang-bank/utils"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load config: %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		fmt.Printf("cannot connect to db: %v", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		fmt.Printf("cannot start server: %v", err)
	}
}
