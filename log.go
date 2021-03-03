package ca_game_api

import (
	"log"
	"net/http"
)

func outputStartLog(r *http.Request) {
	log.Printf("INFO START %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
}

func outputSuccessfulEndLog(r *http.Request) {
	log.Printf("INFO END %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
}