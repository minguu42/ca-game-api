package ca_game_api

import (
	"log"
	"net/http"
)

type GetCharacterListResponse struct {
	Characters []Character `json:"characters"`
}

func GetCharacterList(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(w, r, http.MethodGet) {
		return
	}

	xToken := r.Header.Get("x-token")

	db := Connect()
	defer db.Close()
	characters, err := selectCharacterList(db, xToken, w)
	if err != nil {
		return
	}

	jsonResponse := GetCharacterListResponse{
		Characters: characters,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		return
	}
}

type PutCharacterComposeRequest struct {
	BaseUserCharacterId     int `json:"baseUserCharacterID"`
	MaterialUserCharacterId int `json:"materialUserCharacterID"`
}

type PutCharacterComposeResponse struct {
	UserCharacterId string `json:"userCharacterID"`
	CharacterId     string `json:"characterID"`
	Name            string `json:"name"`
	Level           int    `json:"level"`
}

func PutCharacterCompose(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(w, r, http.MethodPut) {
		return
	}

	xToken := r.Header.Get("x-token")
	var jsonRequest PutCharacterComposeRequest
	if err := decodeRequest(r, &jsonRequest, w); err != nil {
		return
	}
	baseUserCharacterId := jsonRequest.BaseUserCharacterId
	materialUserCharacterId := jsonRequest.MaterialUserCharacterId

	db := Connect()
	defer db.Close()
	userId, err := selectUserId(db, xToken, w)
	if err != nil {
		return
	}
	baseUserId, err := selectUserIdByUserCharacterId(db, baseUserCharacterId, w)
	if err != nil {
		return
	}
	materialUserId, err := selectUserIdByUserCharacterId(db, materialUserCharacterId, w)
	if err != nil {
		return
	}
	if (userId != baseUserId) || (userId != materialUserId) {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 400: User does not own the character")
		return
	}

	tx, newLevel, err := composeCharacter(db, baseUserCharacterId, materialUserCharacterId, w)
	if err != nil {
		if tx != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("ERROR Rollback error:", err)
			}
		}
		return
	}

	jsonResponse, err := createPutCharacterComposeResponse(db, baseUserCharacterId, newLevel, w)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR Rollback error:", err)
		}
		return
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
	log.Println("INFO Commit character compose")
}
