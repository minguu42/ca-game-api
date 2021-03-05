package ca_game_api

import (
	"database/sql"
	"log"
	"net/http"
)

type Character struct {
	UserCharacterId string `json:"userCharacterID"`
	CharacterId     string `json:"characterID"`
	Name            string `json:"name"`
	Level           int    `json:"level"`
}

func selectCharacterName(db *sql.DB, characterId int) (string, error) {
	const selectSql = "SELECT name FROM characters WHERE id = ?"
	var name string
	row := db.QueryRow(selectSql, characterId)
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func countPerRarity(db *sql.DB, rarity int, w http.ResponseWriter) (int, error) {
	const selectSql = "SELECT COUNT(*) FROM characters WHERE rarity = ?"
	var count int
	row := db.QueryRow(selectSql, rarity)
	if err := row.Scan(&count); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return 0, err
	}
	return count, nil
}
