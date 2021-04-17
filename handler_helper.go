package ca_game_api

import (
	"encoding/json"
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

func decodeRequest(r *http.Request, jsonRequest interface{}, w http.ResponseWriter) error {
	if err := json.NewDecoder(r.Body).Decode(jsonRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 400:", err)
		return err
	}
	return nil
}

func encodeResponse(w http.ResponseWriter, jsonResponse interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("INFO Return 500:", err)
		return err
	}
	return nil
}
