package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api/database"
	"github.com/minguu42/ca-game-api/helper"
	"log"
	"net/http"
)

type CharacterListJsonResponse struct {
	Characters []database.Character `json:"characters"`
}

func CharacterListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")
	digestXToken := helper.HashToken(xToken)
	db := database.Connect()
	defer db.Close()

	characters, err := database.GetCharacterList(db, digestXToken)
	if err != nil {
		log.Fatal("get character list error: ", err)
	}
	jsonResponse := CharacterListJsonResponse{
		Characters: characters,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Fatal("json encode error: ", err)
	}
}