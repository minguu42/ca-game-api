package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/minguu42/ca-game-api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/user/create", ca_game_api.PostUser)
	http.HandleFunc("/user/get", ca_game_api.GetUser)
	http.HandleFunc("/user/update", ca_game_api.PutUser)

	http.HandleFunc("/gacha/draw", ca_game_api.GachaDrawHandler)

	http.HandleFunc("/character/list", ca_game_api.CharacterListHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server listen error: ", err)
	}
}
