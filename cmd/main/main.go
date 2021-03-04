package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/minguu42/ca-game-api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/user/create", logging(ca_game_api.PostUser))
	http.HandleFunc("/user/get", logging(ca_game_api.GetUser))
	http.HandleFunc("/user/update", logging(ca_game_api.PutUser))

	http.HandleFunc("/gacha/draw", logging(ca_game_api.PostGachaDraw))

	http.HandleFunc("/character/list", logging(ca_game_api.GetCharacterList))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server listen error: ", err)
	}
}

func logging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO START %v requeest to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
		h(w, r)
		log.Printf("INFO END %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
	}
}
