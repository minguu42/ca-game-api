package ca_game_api

import (
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
	if isRequestMethodInvalid(w, r, "GET") {
		return
	}

	token := r.Header.Get("x-token")

	if _, err := getUserByDigestToken(hash(token)); err != nil {
		w.WriteHeader(401)
		log.Println("ERROR getUserByDigestToken failed:", err)
		return
	}

	userCharacters, err := getUserCharactersByToken(token)
	if err != nil {
		log.Println("ERROR getUserCharactersByToken failed:", err)
		w.WriteHeader(500)
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
	if isRequestMethodInvalid(w, r, "PUT") {
		return
	}

	var reqBody PutCharacterComposeRequest
	if err := decodeRequest(w, r, &reqBody); err != nil {
		return
	}

	baseUserCharacterId := reqBody.BaseUserCharacterId
	materialUserCharacterId := reqBody.MaterialUserCharacterId
	if baseUserCharacterId == materialUserCharacterId {
		log.Println("ERROR cannot composeUserCharacter same character")
		w.WriteHeader(400)
		return
	}

	token := r.Header.Get("x-token")
	user, err := getUserByDigestToken(hash(token))
	if err != nil {
		log.Println("ERROR getUserByDigestToken failed:", err)
		w.WriteHeader(401)
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

	newExperience, err := composeUserCharacter(tx, baseUserCharacter, materialUserCharacter)
	if err != nil {
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
		Level:           calculateLevel(newExperience),
		Experience:      newExperience,
		Power:           calculatePower(newExperience, baseUserCharacter.character.basePower),
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
