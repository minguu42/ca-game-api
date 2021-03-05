package ca_game_api

import (
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
	if err := encodeResponse(w, jsonResponse); err != nil {
		return
	}
}
