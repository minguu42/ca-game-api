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

	results, err, tx := draw(xToken, times, w)
	if err != nil {
		if tx != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("ERROR Rollback error:", err)
			}
		}
		return
	}

	jsonResponse := PostGachaDrawResponse{
		Results: results,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR Rollback error:", err)
		}
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return
	}
	log.Println("INFO Commit gacha result")
}
