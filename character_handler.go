package ca_game_api

import (
	"log"
	"net/http"
)

type GetCharacterListResponse struct {
	Characters []Character2 `json:"characters"`
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
		log.Println("ERROR encodeResponse fail:", err)
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
		log.Println("ERROR decodeRequest fail:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	baseUserCharacterId := jsonRequest.BaseUserCharacterId
	materialUserCharacterId := jsonRequest.MaterialUserCharacterId

	userId, err := selectUserId(xToken)
	if err != nil {
		log.Println("ERROR selectUserId failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	baseUserId, err := selectUserIdByUserCharacterId(baseUserCharacterId)
	if err != nil {
		log.Println("ERROR selectUserIdByUserCharacterId failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	materialUserId, err := selectUserIdByUserCharacterId(materialUserCharacterId)
	if err != nil {
		log.Println("ERROR selectUserIdByUserCharacterId failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if (userId != baseUserId) || (userId != materialUserId) {
		log.Println("ERROR User does not own the character")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx, newLevel, err := composeCharacter(baseUserCharacterId, materialUserCharacterId)
	if err != nil {
		log.Println("ERROR composeCharacter failed:", err)
		if tx != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("ERROR tx.Rollback failed:", err)
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse, err := createPutCharacterComposeResponse(baseUserCharacterId, newLevel)
	if err != nil {
		log.Println("ERROR createPutCharacterComposeResponse failed:", err)
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR tx.Rollback failed:", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		log.Println("ERROR encodeResponse failed:", err)
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR tx.Rollback failed:", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("ERROR tx.Commit failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
