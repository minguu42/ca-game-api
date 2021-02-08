package handlers

import (
	"encoding/json"
	"github.com/minguu42/ca-game-api/database"
	"github.com/minguu42/ca-game-api/helper"
	"log"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}

type Token struct {
	Token string `json:"token"`
}

type UserCreateRequestJson struct {
	Name string `json:"name"`
}

type UserCreateResponseJson struct {
	Token string `json:"token"`
}

func UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// ボディからユーザの取得
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal("user decode error: ", err)
	}

	// トークンの作成
	randomStr, err := helper.GenerateRandomString(22)
	if err != nil {
		log.Fatal("token generate error: ", err)
	}
	token := Token{
		Token: randomStr,
	}

	// トークンをハッシュ化し、データベースに保存する。
	digestToken := helper.HashToken(token.Token)

	db := database.Connect()
	defer db.Close()
	if err := database.InsertUser(db, user.Name, digestToken); err != nil {
		log.Fatal("database create user error: ", err)
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		log.Fatal("json encode error: ", err)
	}
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

	user := User{
		Name: name,
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Fatal("json encode error: ", err)
	}
}

func UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	xToken := r.Header.Get("x-token")
	digestXToken := helper.HashToken(xToken)
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal("user decode error: ", err)
	}

	db := database.Connect()
	defer db.Close()
	id, err := database.GetUserId(db, digestXToken)
	if err != nil {
		log.Fatal("database get user id error: ", err)
	}
	if 	err := database.UpdateUser(db, id, user.Name); err != nil {
		log.Fatal("database user update error: ", err)
	}
}