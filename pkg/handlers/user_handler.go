package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api/pkg/database"
	"github.com/minguu42/ca-game-api/pkg/helper"
	"github.com/minguu42/ca-game-api/pkg/user"
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


	db := database.Connect()
	defer db.Close()
	if err := user.Insert(db, jsonRequest.Name, token); err != nil {
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

	db := database.Connect()
	defer db.Close()
	name, err := user.GetName(db, xToken)
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
	var jsonRequest UserUpdateJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		log.Fatal("user decode error: ", err)
	}

	db := database.Connect()
	defer db.Close()
	if err := user.Update(db, xToken, jsonRequest.Name); err != nil {
		log.Fatal("database user update error: ", err)
	}
}