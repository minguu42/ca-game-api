package ca_game_api

import (
	"database/sql"
	"log"
	"os"
	"time"
)

var db *sql.DB

func OpenDb() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("open db failed:", err)
	}

	isDBReady := false
	failureTimes := 0
	for !isDBReady {
		err := db.Ping()
		if err == nil {
			isDBReady = true
		} else if failureTimes < 2 {
			log.Println("ping db failed. try again")
			time.Sleep(time.Second * 15)
			failureTimes += 1
		} else {
			log.Fatal("ping db failed:", err)
		}
	}
}

func CloseDb() {
	if err := db.Close(); err != nil {
		log.Fatal("close db failed:", err)
	}
}
