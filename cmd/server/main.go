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
		start := time.Now()
		h(w, r)
		end := time.Now()
		log.Printf("INFO %v %v %v in %vμs\n", start.Format("2006/01/02 03:04:05 MST"), r.Method, r.URL, (end.Sub(start)).Microseconds())
	}
}

func main() {
	caGameApi.OpenDb()
	defer caGameApi.CloseDb()

	mux := http.NewServeMux()

	mux.HandleFunc("/user/create", logging(caGameApi.PostUser))
	mux.HandleFunc("/user/get", logging(caGameApi.GetUser))
	mux.HandleFunc("/user/update", logging(caGameApi.PutUser))
	mux.HandleFunc("/user/ranking", logging(caGameApi.GetUserRanking))

	mux.HandleFunc("/gacha/draw", logging(caGameApi.PostGachaDraw))

	mux.HandleFunc("/character/list", logging(caGameApi.GetCharacterList))
	mux.HandleFunc("/character/compose", logging(caGameApi.PutCharacterCompose))

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("s.listenAndServe failed: ", err)
	}
}
