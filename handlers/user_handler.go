package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api/database"
	"github.com/minguu42/ca-game-api/helper"
	"log"
	"net/http"
)

type UserCreateJsonRequest struct {
	Name string `json:"name"`
}

type UserCreateJsonResponse struct {
	Token string `json:"token"`
}

func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var jsonRequest UserCreateJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		log.Fatal("user decode error: ", err)
	}

	token, err := helper.GenerateRandomString(22)
	if err != nil {
		log.Fatal("token generate error: ", err)
	}
	digestToken := helper.HashToken(token)

	db := database.Connect()
	defer db.Close()
	if err := database.InsertUser(db, jsonRequest.Name, digestToken); err != nil {
		log.Fatal("database create user error: ", err)
	}

	jsonResponse := UserCreateJsonResponse{
		Token: token,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Fatal("json encode error: ", err)
	}
}

type UserGetJsonResponse struct {
	Name string `json:"name"`
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")
	digestXToken := helper.HashToken(xToken)

	db := database.Connect()
	defer db.Close()
	name, err := database.GetUserName(db, digestXToken)
	if err != nil {
		log.Fatal("database get user name error: ", err)
	}

	jsonResponse := UserGetJsonResponse{
		Name: name,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Fatal("json encode error: ", err)
	}
}

type UserUpdateJsonRequest struct {
	Name string `json:"name"`
}

func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")
	digestXToken := helper.HashToken(xToken)
	var jsonRequest UserUpdateJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		log.Fatal("user decode error: ", err)
	}

	db := database.Connect()
	defer db.Close()
	id, err := database.GetUserId(db, digestXToken)
	if err != nil {
		log.Fatal("database get user id error: ", err)
	}
	if err := database.UpdateUser(db, id, jsonRequest.Name); err != nil {
		log.Fatal("database user update error: ", err)
	}
}
