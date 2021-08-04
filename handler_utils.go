package ca_game_api

import (
	"encoding/json"
	"log"
	"net/http"
)

func isRequestMethodInvalid(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(405)
		log.Println("ERROR Bad Request Method")
		return true
	}
	return false
}

func decodeRequest(w http.ResponseWriter, r *http.Request, jsonRequest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(jsonRequest); err != nil {
		w.WriteHeader(400)
		log.Println("ERROR decodeRequest failed:", err)
		return err
	}
	return nil
}

func encodeResponse(w http.ResponseWriter, jsonResponse interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonResponse); err != nil {
		w.WriteHeader(500)
		log.Println("ERROR encodeResponse failed:", err)
		return err
	}
	return nil
}
