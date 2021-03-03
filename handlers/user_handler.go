package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api/pkg/database"
	"github.com/minguu42/ca-game-api/pkg/helper"
	"github.com/minguu42/ca-game-api/pkg/user"
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
	outputStartLog(r)
	if isStatusMethodInvalid(w, r, http.MethodPost) {
		return
	}

	var jsonRequest PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json decode error:", err)
		return
	}
	name := jsonRequest.Name

	token, err := helper.GenerateRandomString(22)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Token generate error:", err)
		return
	}

	db := database.Connect()
	defer db.Close()
	if err := user.Insert(db, name, token); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Create user error:", err)
		return
	}

	jsonResponse := PostUserResponse{
		Token: token,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json encode error:", err)
		return
	}
	outputSuccessfulEndLog(r)
}

type GetUserResponse struct {
	Name string `json:"name"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	outputStartLog(r)
	if isStatusMethodInvalid(w, r, http.MethodGet) {
		return
	}

	xToken := r.Header.Get("x-token")

	db := database.Connect()
	defer db.Close()
	name, err := user.GetName(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("ERROR x-token is invalid")
		return
	}

	jsonResponse := GetUserResponse{
		Name: name,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json encode error:", err)
		return
	}
	outputSuccessfulEndLog(r)
}

type PutUserRequest struct {
	Name string `json:"name"`
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	outputStartLog(r)
	if isStatusMethodInvalid(w, r, http.MethodPut) {
		return
	}

	xToken := r.Header.Get("x-token")

	var jsonRequest PutUserRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json decode error: ", err)
		return
	}
	name := jsonRequest.Name

	db := database.Connect()
	defer db.Close()
	user.Update(db, xToken, name, w)

	outputSuccessfulEndLog(r)
}
