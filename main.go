package main

import (
	"database/sql"
	"log"

	"github.com/NghiaLeopard/simple-bank/api"
	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	"github.com/NghiaLeopard/simple-bank/utils"
	_ "github.com/lib/pq"
)


 
func main() {
	config,err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannnot connect to viper",err)
	}

	conn,err := sql.Open(config.DBDrive,config.DBSource)

	if err != nil {
		log.Fatal("Cannnot connect to db",err)
	}

	store := db.NewStore(conn)
	server,err := api.NewServer(config,store)

	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:",err)
	}

}