package main

import (
	"database/sql"
	"log"

	"github.com/NghiaLeopard/simple-bank/api"
	db "github.com/NghiaLeopard/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDrive = "postgres";
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

 
func main() {
	conn,err := sql.Open(dbDrive,dbSource)

	if err != nil {
		log.Fatal("Cannnot connect to db",err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err1 := server.Start(serverAddress)

	if err1 != nil {
		log.Fatal("cannot start server:",err)
	}

}