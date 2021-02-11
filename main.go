package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/minguu42/ca-game-api/database"
	"github.com/minguu42/ca-game-api/handlers"
	"log"
	"net/http"
)

func main() {
	db := database.Connect()
	defer db.Close()

	// 通信確認用。
	if err := db.Ping(); err != nil {
		log.Fatal("database ping error: ", err)
	}

	http.HandleFunc("/user/create", handlers.UserCreateHandler)
	http.HandleFunc("/user/get", handlers.UserGetHandler)
	http.HandleFunc("/user/update", handlers.UserUpdateHandler)

	http.HandleFunc("/gacha/draw", handlers.GachaDrawHandler)

	http.HandleFunc("/character/list", handlers.CharacterListHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("server error: ", err)
	}
}
