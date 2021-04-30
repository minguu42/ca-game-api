package ca_game_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func isStatusMethodInvalid(r *http.Request, method string) bool {
	if r.Method != method {
		return true
	}
	return false
}

func decodeRequest(r *http.Request, jsonRequest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(jsonRequest); err != nil {
		return fmt.Errorf("decode fail: %w", err)
	}
	return nil
}

func encodeResponse(w http.ResponseWriter, jsonResponse interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonResponse); err != nil {
		return fmt.Errorf("encode fail: %w", err)
	}
	return nil
}
