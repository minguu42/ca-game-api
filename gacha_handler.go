package ca_game_api

import (
	"log"
	"net/http"
)

type PostGachaDrawRequest struct {
	Times int `json:"times"`
}

type ResultJson struct {
	CharacterId int    `json:"characterID"`
	Name        string `json:"name"`
}

type PostGachaDrawResponse struct {
	Results []ResultJson `json:"results"`
}

func PostGachaDraw(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, http.MethodPost) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")
	var reqBody PostGachaDrawRequest
	if err := decodeRequest(r, &reqBody); err != nil {
		log.Println("ERROR decodeRequest fail:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	times := reqBody.Times

	if times <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 403: Times is 0 or negative number")
		return
	}

	results, err, tx := draw(xToken, times)
	if err != nil {
		if tx != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("ERROR Rollback error:", err)
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody := PostGachaDrawResponse{
		Results: results,
	}
	if err := encodeResponse(w, respBody); err != nil {
		log.Println("ERROR encodeResponse fail:", err)
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR Rollback error:", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return
	}
}
