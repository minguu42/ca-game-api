package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api/pkg/database"
	"github.com/minguu42/ca-game-api/pkg/gacha"
	"github.com/minguu42/ca-game-api/pkg/user"
	"log"
	"net/http"
)

type GachaDrawJsonRequest struct {
	Times int `json:"times"`
}

type GachaDrawJsonResponse struct {
	Results []gacha.Result `json:"results"`
}

func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")
	var jsonRequest GachaDrawJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		log.Fatal("user decode error: ", err)
	}
	times := jsonRequest.Times

	db := database.Connect()
	defer db.Close()
	userId, err := user.GetId(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	results, err := gacha.Draw(db, userId, times)
	if err != nil {
		log.Fatal("gacha draw error")
	}

	jsonResponse := GachaDrawJsonResponse{
		Results: results,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Fatal("json encode error: ", err)
	}
}
