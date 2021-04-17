package ca_game_api

import (
	"database/sql"
	"log"
	"os"
)

func Connect() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("database connection error", err)
	}
	return db
}
