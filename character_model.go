package ca_game_api

import "database/sql"

func selectCharacterName(db *sql.DB, characterId int) (string, error) {
	const selectSql = "SELECT name FROM characters WHERE id = ?"
	var name string
	row := db.QueryRow(selectSql, characterId)
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func countPerRarity(db *sql.DB, rarity int) (int, error) {
	const countSql = "SELECT COUNT(*) FROM characters WHERE rarity = ?"
	var count int
	row := db.QueryRow(countSql, rarity)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
