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

	var jsonRequest PostGachaDrawRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 403:", err)
		return
	}
	times := jsonRequest.Times
	if times <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 403: Times is 0 or negative number")
		return
	}

	db := Connect()
	defer db.Close()
	userId, err := selectUserId(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("ERROR Ruturn 401: x-token is invalid")
		return
	}

	results, err := Draw(db, userId, times)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return
	}

	jsonResponse := PostGachaDrawResponse{
		Results: results,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("INFO Return 500:", err)
	}
}
