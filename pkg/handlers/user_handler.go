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
	log.Printf("INFO START %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("ERROR Status method is not allowed")
		return
	}

	var jsonRequest UserCreateJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		log.Println("ERROR Json decode error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	name := jsonRequest.Name
	log.Println("INFO Get user name - Success")

	token, err := helper.GenerateRandomString(22)
	if err != nil {
		log.Println("ERROR Token generate error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("INFO Generate token - Success")

	db := database.Connect()
	defer db.Close()
	if err := user.Insert(db, name, token); err != nil {
		log.Println("ERROR Create user error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("INFO Create user - Success")

	jsonResponse := UserCreateJsonResponse{
		Token: token,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Println("ERROR Json encode error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("INFO END %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
}

type UserGetJsonResponse struct {
	Name string `json:"name"`
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO START %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("ERROR Status method is not allowed")
		return
	}

	xToken := r.Header.Get("x-token")
	log.Println("INFO Get x-token - Success")

	db := database.Connect()
	defer db.Close()
	name, err := user.GetName(db, xToken)
	if err != nil {
		log.Println("ERROR x-token is invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("INFO Get user name - Success")

	jsonResponse := UserGetJsonResponse{
		Name: name,
	}
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Println("ERROR Json encode error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("INFO END %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
}

type UserUpdateJsonRequest struct {
	Name string `json:"name"`
}

func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO START %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Println("ERROR Status method is not allowed")
		return
	}

	xToken := r.Header.Get("x-token")
	log.Println("INFO Get x-token - Success")

	var jsonRequest UserUpdateJsonRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
		log.Println("ERROR Json decode error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	name := jsonRequest.Name
	log.Println("INFO Get user name - Success")

	db := database.Connect()
	defer db.Close()
	if err := user.Update(db, xToken, name); err != nil {
		log.Println("ERROR Update user error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("INFO Update user - Success")

	log.Printf("INFO END %v request to %v came from %v", r.Method, r.URL, r.Header.Get("User-Agent"))
}
