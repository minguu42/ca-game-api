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
	if isStatusMethodInvalid(r, "POST") {
		w.WriteHeader(405)
		return
	}

	token := r.Header.Get("x-token")
	var reqBody PostGachaDrawRequest
	if err := decodeRequest(r, &reqBody); err != nil {
		log.Println("ERROR decodeRequest failed:", err)
		w.WriteHeader(400)
		return
	}
	times := reqBody.Times

	if times <= 0 {
		w.WriteHeader(400)
		log.Println("ERROR Times should be positive number")
		return
	}
	user, err := getUserByToken(token)
	if err != nil {
		log.Println("ERROR getUserByToken failed:", err)
		w.WriteHeader(403)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("ERROR db.Begin failed:", err)
		w.WriteHeader(500)
		return
	}

	results, err := draw(tx, user.id, times)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR tx.Rollback failed:", err)
		}
		w.WriteHeader(500)
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
