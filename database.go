package ca_game_api

import (
	"database/sql"
	"log"
	"os"
)

var (
	DRIVER = os.Getenv("DRIVER")
	DATASOURCE = os.Getenv("DATASOURCE")
)

func Connect() *sql.DB {
	db, err := sql.Open(DRIVER, DATASOURCE)
	if err != nil {
		log.Fatal("database connection error", err)
	}
	return db
}