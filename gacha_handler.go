package ca_game_api

import (
	"encoding/json"
	"log"
	"net/http"
)

type PostGachaDrawRequest struct {
	Times int `json:"times"`
}

type PostGachaDrawResponse struct {
	Results []Result `json:"results"`
}

func PostGachaDraw(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(w, r, http.MethodPost) {
		return
	}

	xToken := r.Header.Get("x-token")
	log.Println("INFO Get x-token - Success")

	var jsonRequest PostGachaDrawRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json decode error:", err)
		return
	}
	times := jsonRequest.Times
	log.Println("INFO Get gacha times - Success")

	db := Connect()
	defer db.Close()
	userId, err := selectUserId(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("ERROR x-token is invalid")
		return
	}
	log.Println("INFO Get userId - Success")

	results, err := Draw(db, userId, times)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Draw gacha error:", err)
		return
	}
	log.Println("INFO Draw gacha - Success")

	jsonResponse := PostGachaDrawResponse{
		Results: results,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("INFO Json encode error:", err)
		return
	}
}
