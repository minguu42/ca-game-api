package ca_game_api

import (
	"log"
	"net/http"
)

func Log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO START %v requeest to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
		h(w, r)
		log.Printf("INFO END %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
	}
}
