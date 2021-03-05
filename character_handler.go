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
	BaseUserCharacterId int `json:"baseCharacterID"`
	MaterialUserCharacterId int `json:"materialUserCharacterID"`
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
	log.Println(xToken, baseUserCharacterId, materialUserCharacterId)
}