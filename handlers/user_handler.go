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
	outputStartLog(r)
	if isStatusMethodInvalid(w, r, http.MethodPost) {
		return
	}

	var jsonRequest UserCreateJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json decode error:", err)
		return
	}
	name := jsonRequest.Name
	log.Println("INFO Get user name - Success")

	token, err := helper.GenerateRandomString(22)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Token generate error:", err)
		return
	}
	log.Println("INFO Generate token - Success")

	db := database.Connect()
	defer db.Close()
	if err := user.Insert(db, name, token); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Create user error:", err)
		return
	}
	log.Println("INFO Create user - Success")

	jsonResponse := UserCreateJsonResponse{
		Token: token,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json encode error:", err)
		return
	}
	outputSuccessfulEndLog(r)
}

type UserGetJsonResponse struct {
	Name string `json:"name"`
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	outputStartLog(r)
	if isStatusMethodInvalid(w, r, http.MethodGet) {
		return
	}

	xToken := r.Header.Get("x-token")
	log.Println("INFO Get x-token - Success")

	db := database.Connect()
	defer db.Close()
	name, err := user.GetName(db, xToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("ERROR x-token is invalid")
		return
	}
	log.Println("INFO Get user name - Success")

	jsonResponse := UserGetJsonResponse{
		Name: name,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json encode error:", err)
		return
	}
	outputSuccessfulEndLog(r)
}

type UserUpdateJsonRequest struct {
	Name string `json:"name"`
}

func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	outputStartLog(r)
	if isStatusMethodInvalid(w, r, http.MethodPut) {
		return
	}

	xToken := r.Header.Get("x-token")
	log.Println("INFO Get x-token - Success")

	var jsonRequest UserUpdateJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Json decode error: ", err)
		return
	}
	name := jsonRequest.Name
	log.Println("INFO Get user name - Success")

	db := database.Connect()
	defer db.Close()
	user.Update(db, xToken, name, w)
	log.Println("INFO END Update user")

	outputSuccessfulEndLog(r)
}
