package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	caGameApi "github.com/minguu42/ca-game-api"
)

func main() {
	http.HandleFunc("/user/create", measure(logging(caGameApi.PostUser)))
	http.HandleFunc("/user/get", measure(logging(caGameApi.GetUser)))
	http.HandleFunc("/user/update", measure(logging(caGameApi.PutUser)))

	http.HandleFunc("/gacha/draw", measure(logging(caGameApi.PostGachaDraw)))

	http.HandleFunc("/ranking/user", measure(logging(caGameApi.GetRankingUser)))

	http.HandleFunc("/character/list", measure(logging(caGameApi.GetCharacterList)))
	http.HandleFunc("/character/compose", measure(logging(caGameApi.PutCharacterCompose)))

	if err := http.ListenAndServe(":" + os.Getenv("PORT"), nil); err != nil {
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
