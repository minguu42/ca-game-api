package ca_game_api

import (
	"log"
	"net/http"
)

type PostUserRequest struct {
	Name string `json:"name"`
}

type PostUserResponse struct {
	Token string `json:"token"`
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, http.MethodPost) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var jsonRequest PostUserRequest
	if err := decodeRequest(r, &jsonRequest); err != nil {
		log.Println("ERROR decodeRequest failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := generateRandomString(22)
	if err != nil {
		log.Println("ERROR generateRandomString failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var user User
	user.name = jsonRequest.Name
	user.digestToken = hash(token)

	if err := insertUser(user); err != nil {
		log.Println("ERROR insertUser failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonResponse := PostUserResponse{
		Token: token,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		log.Println("ERROR encodeResponse failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type GetUserResponse struct {
	Name string `json:"name"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, http.MethodGet) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("x-token")

	user, err := selectUserByToken(token)
	if err != nil {
		log.Println("ERROR selectUserByToken failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonResponse := GetUserResponse{
		Name: user.name,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		log.Println("ERROR encodeResponse failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type PutUserRequest struct {
	Name string `json:"name"`
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, http.MethodPut) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("x-token")
	var jsonRequest PutUserRequest
	if err := decodeRequest(r, &jsonRequest); err != nil {
		log.Println("ERROR decodeRequest failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	user.name = jsonRequest.Name
	user.digestToken = hash(token)

	if err := updateUser(user); err != nil {
		log.Println("ERROR updateUser failed:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
