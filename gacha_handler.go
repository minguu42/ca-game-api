package ca_game_api

import (
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
	if err := decodeRequest(r, &jsonRequest, w); err != nil {
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

	results, err := draw(db, xToken, times, w)
	if err != nil {
		return
	}

	jsonResponse := PostGachaDrawResponse{
		Results: results,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		return
	}
}
