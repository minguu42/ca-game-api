package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	caGameApi "github.com/minguu42/ca-game-api"
)

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

func main() {
	caGameApi.OpenDb()
	defer caGameApi.CloseDb()

	mux := http.NewServeMux()

	mux.HandleFunc("/user/create", measure(logging(caGameApi.PostUser)))
	mux.HandleFunc("/user/get", measure(logging(caGameApi.GetUser)))
	mux.HandleFunc("/user/update", measure(logging(caGameApi.PutUser)))
	mux.HandleFunc("/user/ranking", measure(logging(caGameApi.GetUserRanking)))

	mux.HandleFunc("/gacha/draw", measure(logging(caGameApi.PostGachaDraw)))

	mux.HandleFunc("/character/list", measure(logging(caGameApi.GetCharacterList)))
	mux.HandleFunc("/character/compose", measure(logging(caGameApi.PutCharacterCompose)))

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("Server listen error: ", err)
	}
}
