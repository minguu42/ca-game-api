package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	ca_game_api "github.com/minguu42/ca-game-api"
)

func main() {
	http.HandleFunc("/user/create", measure(logging(ca_game_api.PostUser)))
	http.HandleFunc("/user/get", measure(logging(ca_game_api.GetUser)))
	http.HandleFunc("/user/update", measure(logging(ca_game_api.PutUser)))

	http.HandleFunc("/gacha/draw", measure(logging(ca_game_api.PostGachaDraw)))

	http.HandleFunc("/character/list", measure(logging(ca_game_api.GetCharacterList)))
	http.HandleFunc("/character/compose", measure(logging(ca_game_api.PutCharacterCompose)))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server listen error: ", err)
	}
}

func logging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO START %v requeest to %v\n", r.Method, r.URL)
		h(w, r)
		log.Printf("INFO END %v request to %v\n", r.Method, r.URL)
	}
}

func measure(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		end := time.Now()
		log.Printf("INFO Response time: %v seconds\n", (end.Sub(start)).Seconds())
	}
}
