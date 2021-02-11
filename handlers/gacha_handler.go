package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api/database"
	"github.com/minguu42/ca-game-api/helper"
	"log"
	"net/http"
	"strconv"
)

type GachaResult struct {
	CharacterId string `json:"characterID"`
	Name        string `json:"name"`
}

type GachaDrawJsonRequest struct {
	Times int `json:"times"`
}

type GachaDrawJsonResponse struct {
	Results *[]GachaResult `json:"results"`
}

func GachaDrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")
	digestXToken := helper.HashToken(xToken)
	var jsonRequest GachaDrawJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		log.Fatal("user decode error: ", err)
	}
	times := jsonRequest.Times

	db := database.Connect()
	defer db.Close()
	rarity3CharacterNum, err := database.CountCharacterPerRarity(db, 3)
	if err != nil {
		log.Fatal("database get count error: ", err)
	}
	rarity4CharacterNum, err := database.CountCharacterPerRarity(db, 3)
	if err != nil {
		log.Fatal("database get count error: ", err)
	}
	rarity5CharacterNum, err := database.CountCharacterPerRarity(db, 3)
	if err != nil {
		log.Fatal("database get count error: ", err)
	}
	var selectedCharacterId int
	switch selectedRarity := helper.SelectRarity(); selectedRarity {
	case 3:
		selectedCharacterId = helper.SelectCharacterId(rarity3CharacterNum) + 30000000
	case 4:
		selectedCharacterId = helper.SelectCharacterId(rarity4CharacterNum) + 40000000
	case 5:
		selectedCharacterId = helper.SelectCharacterId(rarity5CharacterNum) + 50000000
	}
	name, err := database.GetCharacterName(db, selectedCharacterId)
	if err != nil {
		log.Fatal("database get character name error: ", err)
	}
	jsonResponse := GachaResult{
		CharacterId: strconv.Itoa(selectedCharacterId),
		Name: name,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Fatal("json encode error: ", err)
	}
	println(digestXToken, times)
}
