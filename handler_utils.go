package ca_game_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func isRequestMethodInvalid(r *http.Request, method string) bool {
	return r.Method != method
}

func decodeRequest(r *http.Request, jsonRequest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(jsonRequest); err != nil {
		return fmt.Errorf("decoder.Decode failed: %w", err)
	}
	return nil
}

func encodeResponse(w http.ResponseWriter, jsonResponse interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonResponse); err != nil {
		return fmt.Errorf("encoder.Encode failed: %w", err)
	}
	return nil
}
