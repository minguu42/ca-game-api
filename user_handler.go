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
	if isRequestMethodInvalid(r, "POST") {
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

	if err := insertUser(reqBody.Name, hash(token)); err != nil {
		fmt.Println("ERROR insertUser")
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
	if isRequestMethodInvalid(r, "GET") {
		w.WriteHeader(405)
		return
	}

	token := r.Header.Get("x-token")

	user, err := getUserByDigestToken(hash(token))
	if err != nil {
		log.Println("ERROR getUserByDigestToken failed:", err)
		w.WriteHeader(403)
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
	if isRequestMethodInvalid(r, "PUT") {
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

	user := User{
		name:        reqBody.Name,
		digestToken: hash(token),
	}
	if err := updateUser(user); err != nil {
		log.Println("ERROR updateUser failed:", err)
		w.WriteHeader(400)
		return
	}
}

type UserRankingJson struct {
	Name     string `json:"name"`
	SumPower int    `json:"sumPower"`
}

type GetUserRankingResponse struct {
	Users []UserRankingJson `json:"users"`
}

func GetUserRanking(w http.ResponseWriter, r *http.Request) {
	if isRequestMethodInvalid(r, "GET") {
		w.WriteHeader(405)
		return
	}

	token := r.Header.Get("x-token")
	if _, err := getUserByDigestToken(hash(token)); err != nil {
		fmt.Println("ERROR getUserByDigestToken failed:", err)
		w.WriteHeader(403)
		return
	}

	rankings, err := selectUserRanking()
	if err != nil {
		log.Println("ERROR selectUserRanking failed:", err)
		w.WriteHeader(500)
		return
	}

	jsonResponse := GetUserRankingResponse{
		Users: rankings,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		log.Println("ERROR encodeResponse failed:", err)
		w.WriteHeader(500)
		return
	}
}
