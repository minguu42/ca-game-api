package ca_game_api

import "fmt"

type Character struct {
	id        int
	name      string
	rarity    int
	basePower int
	calorie   int
}


func getCharacterById(id int) (Character, error) {
	const query = "SELECT id, name, rarity, base_power, calorie FROM characters WHERE id = $1"
	row := db.QueryRow(query, id)
	var character Character
	if err := row.Scan(&character.id, &character.name, &character.rarity, &character.basePower, &character.calorie); err != nil {
		return Character{}, fmt.Errorf("row.Scan failed: %w", err)
	}
	return character, nil
}

func countCharactersByRarity(rarity int) (int, error) {
	const query = "SELECT COUNT(*) FROM characters WHERE rarity = $1"
	var count int
	row := db.QueryRow(query, rarity)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("row.Scan faild: %w", err)
	}
	return count, nil
}