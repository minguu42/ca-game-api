package ca_game_api

import (
	"log"
	"net/http"
)

type GetCharacterListResponse struct {
	Characters []Character `json:"characters"`
}

func GetCharacterList(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, http.MethodGet) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")

	characters, err := selectCharacterList(xToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := GetCharacterListResponse{
		Characters: characters,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	if isStatusMethodInvalid(r, http.MethodPut) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")
	var jsonRequest PutCharacterComposeRequest
	if err := decodeRequest(r, &jsonRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	baseUserCharacterId := jsonRequest.BaseUserCharacterId
	materialUserCharacterId := jsonRequest.MaterialUserCharacterId

	userId, err := selectUserId(xToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	baseUserId, err := selectUserIdByUserCharacterId(baseUserCharacterId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	materialUserId, err := selectUserIdByUserCharacterId(materialUserCharacterId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if (userId != baseUserId) || (userId != materialUserId) {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 400: User does not own the character")
		return
	}

	tx, newLevel, err := composeCharacter(baseUserCharacterId, materialUserCharacterId)
	if err != nil {
		if tx != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("ERROR Rollback error:", err)
			}
		}
		return
	}

	jsonResponse, err := createPutCharacterComposeResponse(baseUserCharacterId, newLevel)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return
	}
	log.Println("INFO Commit character compose")
}
