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
	characters, err := selectCharacterList(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 403:", err)
		return
	}

	jsonResponse := GetCharacterListResponse{
		Characters: characters,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return
	}
}
