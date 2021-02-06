package main

import (
	"ca-game-api/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.UserHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("server error: ", err)
	}
}
