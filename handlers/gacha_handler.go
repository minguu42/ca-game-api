package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api"
	"log"
	"net/http"
)

type GachaDrawJsonRequest struct {
	Times int `json:"times"`
}

type GachaDrawJsonResponse struct {
	Results []ca_game_api.Result `json:"results"`
}

func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	outputStartLog(r)
	if isStatusMethodInvalid(w, r, http.MethodPost) {
		return
	}

	xToken := r.Header.Get("x-token")
	log.Println("INFO Get x-token - Success")

	var jsonRequest GachaDrawJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json decode error:", err)
		return
	}
	times := jsonRequest.Times
	log.Println("INFO Get gacha times - Success")

	db := ca_game_api.Connect()
	defer db.Close()
	userId, err := ca_game_api.GetUserId(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("ERROR x-token is invalid")
		return
	}
	log.Println("INFO Get userId - Success")

	results, err := ca_game_api.Draw(db, userId, times)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Draw gacha error:", err)
		return
	}
	log.Println("INFO Draw gacha - Success")

	jsonResponse := GachaDrawJsonResponse{
		Results: results,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("INFO Json encode error:", err)
		return
	}
	outputSuccessfulEndLog(r)
}
