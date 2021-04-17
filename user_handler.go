package ca_game_api

import (
	"net/http"
)

type PostUserRequest struct {
	Name string `json:"name"`
}

type PostUserResponse struct {
	Token string `json:"token"`
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(w, r, http.MethodPost) {
		return
	}

	var jsonRequest PostUserRequest
	if err := decodeRequest(r, &jsonRequest, w); err != nil {
		return
	}
	name := jsonRequest.Name

	token, err := generateRandomString(22, w)
	if err != nil {
		return
	}

	if err := insertUser(name, token, w); err != nil {
		return
	}

	jsonResponse := PostUserResponse{
		Token: token,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		return
	}
}

type GetUserResponse struct {
	Name string `json:"name"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(w, r, http.MethodGet) {
		return
	}

	xToken := r.Header.Get("x-token")

	name, err := selectUserName(xToken, w)
	if err != nil {
		return
	}

	jsonResponse := GetUserResponse{
		Name: name,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		return
	}
}

type PutUserRequest struct {
	Name string `json:"name"`
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(w, r, http.MethodPut) {
		return
	}

	xToken := r.Header.Get("x-token")

	var jsonRequest PutUserRequest
	if err := decodeRequest(r, &jsonRequest, w); err != nil {
		return
	}
	name := jsonRequest.Name

	if err := updateUser(xToken, name, w); err != nil {
		return
	}
}
