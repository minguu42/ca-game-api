package ca_game_api

import (
	"encoding/json"
	"log"
	"net/http"
)

type GetCharacterListResponse struct {
	Characters []Character `json:"characters"`
}

func GetCharacterList(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(w, r, http.MethodGet) {
		return
	}

	xToken := r.Header.Get("x-token")

	db := Connect()
	defer db.Close()
	characters, err := selectCharacterList(db, xToken, w)
	if err != nil {
		return
	}

	jsonResponse := GetCharacterListResponse{
		Characters: characters,
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("INFO Return 500:", err)
		return
	}
}
