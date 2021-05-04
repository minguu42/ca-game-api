package ca_game_api

import (
	"fmt"
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
	if isStatusMethodInvalid(r, "POST") {
		w.WriteHeader(405)
		return
	}

	var reqBody PostUserRequest
	if err := decodeRequest(r, &reqBody); err != nil {
		log.Println("ERROR decodeRequest failed:", err)
		w.WriteHeader(400)
		return
	}

	token, err := generateRandomString(22)
	if err != nil {
		log.Println("ERROR generateRandomString failed:", err)
		w.WriteHeader(500)
		return
	}

	var user User
	user.name = reqBody.Name
	user.digestToken = hash(token)
	if err := user.insert(); err != nil {
		log.Println("ERROR user.insert failed:", err)
		w.WriteHeader(400)
		return
	}

	respBody := PostUserResponse{
		Token: token,
	}
	if err := encodeResponse(w, respBody); err != nil {
		log.Println("ERROR encodeResponse failed:", err)
		w.WriteHeader(500)
		return
	}
}

type GetUserResponse struct {
	Name string `json:"name"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, "GET") {
		w.WriteHeader(405)
		return
	}

	token := r.Header.Get("x-token")

	user, err := getUserByToken(token)
	if err != nil {
		log.Println("ERROR getUserByToken failed:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	respBody := GetUserResponse{
		Name: user.name,
	}
	if err := encodeResponse(w, respBody); err != nil {
		log.Println("ERROR encodeResponse failed:", err)
		w.WriteHeader(500)
		return
	}
}

type PutUserRequest struct {
	Name string `json:"name"`
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, "PUT") {
		w.WriteHeader(405)
		return
	}

	token := r.Header.Get("x-token")
	var reqBody PutUserRequest
	if err := decodeRequest(r, &reqBody); err != nil {
		log.Println("ERROR decodeRequest failed:", err)
		w.WriteHeader(400)
		return
	}

	var user User
	user.name = reqBody.Name
	user.digestToken = hash(token)

	if err := user.update(); err != nil {
		log.Println("ERROR updateUser failed:", err)
		w.WriteHeader(400)
		return
	}
}

type UserJson struct {
	Name     string `json:"name"`
	SumPower string `json:"sumPower"`
}

type GetUserRankingResponse struct {
	Users []UserJson `json:"users"`
}

func GetUserRanking(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, "GET") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("x-token")
	if _, err := getUserByToken(token); err != nil {
		fmt.Println("getUserByToken failed:", err)
		w.WriteHeader(403)
		return
	}

	userRankings, err := selectUserRanking()
	if err != nil {
		log.Println("ERROR selectUserRanking error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := GetUserRankingResponse{
		Users: userRankings,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		log.Println("ERROR encodeResponse fail:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
