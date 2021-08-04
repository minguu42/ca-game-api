package ca_game_api

import (
	"fmt"
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
	if isRequestMethodInvalid(w, r, "POST") {
		return
	}

	token := r.Header.Get("x-token")
	var reqBody PostGachaDrawRequest
	if err := decodeRequest(w, r, &reqBody); err != nil {
		return
	}
	times := reqBody.Times

	if times <= 0 {
		w.WriteHeader(400)
		log.Println("ERROR Times should be positive number")
		return
	}
	user, err := getUserByDigestToken(hash(token))
	if err != nil {
		log.Println("ERROR getUserByDigestToken failed:", err)
		w.WriteHeader(403)
		return
	}

	results, err := decideGachaResults(user, times)
	if err != nil {
		fmt.Println("ERROR decideGachaResults failed:", err)
		w.WriteHeader(500)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("ERROR db.Begin failed:", err)
		w.WriteHeader(500)
		return
	}
	if err := insertGachaResults(tx, results); err != nil {
		log.Println("ERROR insertGachaResults failed:", err)
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR tx.Rollback failed:", err)
		}
		w.WriteHeader(500)
		return
	}

	resultsJson := make([]ResultJson, 0, len(results))
	for _, result := range results {
		resultJson := ResultJson{
			CharacterId: result.character.id,
			Name:        result.character.name,
		}
		resultsJson = append(resultsJson, resultJson)
	}
	respBody := PostGachaDrawResponse{
		Results: resultsJson,
	}
	if err := encodeResponse(w, respBody); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR tx.Rollback failed:", err)
		}
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("ERROR tx.Commit failed:", err)
		w.WriteHeader(500)
		return
	}
}
