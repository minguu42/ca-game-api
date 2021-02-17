package handlers

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

func isStatusMethodInvalid(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("ERROR Status method is not allowed")
		return true
	}
	return false
}
