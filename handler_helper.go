package ca_game_api

import (
	"encoding/json"
	"log"
	"net/http"
)

func isStatusMethodInvalid(r *http.Request, method string) bool {
	if r.Method != method {
		log.Println("ERROR Status method is not allowed")
		return true
	}
	return false
}

func decodeRequest(r *http.Request, jsonRequest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(jsonRequest); err != nil {
		log.Println("ERROR decodeRequest error:", err)
		return err
	}
	return nil
}

func encodeResponse(w http.ResponseWriter, jsonResponse interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonResponse); err != nil {
		log.Println("ERROR encodeResponse error:", err)
		return err
	}
	return nil
}
