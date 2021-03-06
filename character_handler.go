package ca_game_api

import (
	"fmt"
	"log"
	"net/http"
)

type CharacterJson struct {
	UserCharacterId int    `json:"userCharacterID"`
	CharacterId     int    `json:"characterID"`
	Name            string `json:"name"`
	Level           int    `json:"level"`
	Experience      int    `json:"experience"`
	Power           int    `json:"power"`
}

type GetCharacterListResponse struct {
	Characters []CharacterJson `json:"characters"`
}

func GetCharacterList(w http.ResponseWriter, r *http.Request) {
	if isRequestMethodInvalid(r, "GET") {
		w.WriteHeader(405)
		return
	}

	xToken := r.Header.Get("x-token")

	userCharacters, err := getUserCharactersByToken(xToken)
	if err != nil {
		log.Println("ERROR getUserCharactersByToken failed:", err)
		w.WriteHeader(400)
		return
	}

	charactersJson := make([]CharacterJson, 0, len(userCharacters))
	for _, userOwnCharacter := range userCharacters {
		characterJson := CharacterJson{
			UserCharacterId: userOwnCharacter.id,
			CharacterId:     userOwnCharacter.character.id,
			Name:            userOwnCharacter.character.name,
			Level:           calculateLevel(userOwnCharacter.experience),
			Experience:      userOwnCharacter.experience,
			Power:           calculatePower(userOwnCharacter.experience, userOwnCharacter.character.basePower),
		}
		charactersJson = append(charactersJson, characterJson)
	}
	respBody := GetCharacterListResponse{
		Characters: charactersJson,
	}
	if err := encodeResponse(w, respBody); err != nil {
		log.Println("ERROR encodeResponse failed:", err)
		w.WriteHeader(500)
		return
	}
}

type PutCharacterComposeRequest struct {
	BaseUserCharacterId     int `json:"baseUserCharacterID"`
	MaterialUserCharacterId int `json:"materialUserCharacterID"`
}

type PutCharacterComposeResponse struct {
	UserCharacterId int    `json:"userCharacterID"`
	CharacterId     int    `json:"characterID"`
	Name            string `json:"name"`
	Level           int    `json:"level"`
	Experience      int    `json:"experience"`
	Power           int    `json:"power"`
}

func PutCharacterCompose(w http.ResponseWriter, r *http.Request) {
	if isRequestMethodInvalid(r, "PUT") {
		w.WriteHeader(405)
		return
	}

	var reqBody PutCharacterComposeRequest
	if err := decodeRequest(r, &reqBody); err != nil {
		log.Println("ERROR decodeRequest failed:", err)
		w.WriteHeader(400)
		return
	}
	token := r.Header.Get("x-token")
	baseUserCharacterId := reqBody.BaseUserCharacterId
	materialUserCharacterId := reqBody.MaterialUserCharacterId
	if baseUserCharacterId == materialUserCharacterId {
		log.Println("ERROR cannot composeUserCharacter same character")
		w.WriteHeader(400)
		return
	}

	user, err := getUserByDigestToken(hash(token))
	if err != nil {
		fmt.Println("ERROR getUserByDigestToken failed:", err)
		w.WriteHeader(403)
		return
	}
	baseUserCharacter, err := getUserCharacterById(baseUserCharacterId)
	if err != nil {
		log.Println("ERROR getUserCharacterById failed:", err)
		w.WriteHeader(400)
		return
	}
	materialUserCharacter, err := getUserCharacterById(materialUserCharacterId)
	if err != nil {
		log.Println("ERROR getUserCharacterById failed:", err)
		w.WriteHeader(400)
		return
	}
	if (user.id != baseUserCharacter.user.id) || (user.id != materialUserCharacter.user.id) {
		log.Println("ERROR User does not own the character")
		w.WriteHeader(403)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("ERROR db.Begin failed:", err)
		w.WriteHeader(500)
		return
	}

	if err := composeUserCharacter(tx, baseUserCharacter, materialUserCharacter); err != nil {
		log.Println("baseUserCharacter.composeUserCharacter failed:", err)
		if err := tx.Rollback(); err != nil {
			log.Println("ERROR tx.Rollback failed:", err)
		}
		w.WriteHeader(500)
		return
	}

	respBody := PutCharacterComposeResponse{
		UserCharacterId: baseUserCharacter.id,
		CharacterId:     baseUserCharacter.character.id,
		Name:            baseUserCharacter.character.name,
		Level:           calculateLevel(baseUserCharacter.experience),
		Experience:      baseUserCharacter.experience,
		Power:           calculatePower(baseUserCharacter.experience, baseUserCharacter.character.basePower),
	}
	if err := encodeResponse(w, respBody); err != nil {
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
