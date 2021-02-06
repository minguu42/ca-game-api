package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Name string `json:"name"`
}

// 動作確認用。のちに削除する。
func UserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{"John"}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Fatal("json encode error: ", err)
	}
}