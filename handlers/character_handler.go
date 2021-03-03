package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api/pkg/database"
	"github.com/minguu42/ca-game-api/pkg/user"
	"log"
	"net/http"
)

type CharacterListJsonResponse struct {
	Characters []user.Character `json:"characters"`
}

func CharacterListHandler(w http.ResponseWriter, r *http.Request) {
	outputStartLog(r)
	if isStatusMethodInvalid(w, r, http.MethodGet) {
		return
	}

	xToken := r.Header.Get("x-token")
	log.Println("INFO Get x-token - Success")

	db := database.Connect()
	defer db.Close()
	characters, err := user.GetCharacterList(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Get character list error:", err)
		return
	}
	log.Println("INFO Get character list - Success")

	jsonResponse := CharacterListJsonResponse{
		Characters: characters,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Println("ERROR Json encode error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	outputSuccessfulEndLog(r)
}