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
	log.Printf("INFO START %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("ERROR Status method is not allowed")
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

	db := database.Connect()
	defer db.Close()
	userId, err := user.GetId(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("ERROR x-token is invalid")
		return
	}
	log.Println("INFO Get userId - Success")

	results, err := gacha.Draw(db, userId, times)
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
	log.Printf("INFO END %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
}
