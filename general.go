package ca_game_api

import (
	"log"
	"net/http"
)

func isStatusMethodInvalid(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("ERROR Status method is not allowed")
		return true
	}
	return false
}
