package ca_game_api

import (
	"database/sql"
	"fmt"
	"math"
)

type Character struct {
	UserCharacterId string `json:"userCharacterID"`
	CharacterId     string `json:"characterID"`
	Name            string `json:"name"`
	Level           int    `json:"level"`
}

func selectCharacterName(characterId int) (string, error) {
	const selectSql = "SELECT name FROM characters WHERE id = $1"
	var name string
	row := db.QueryRow(selectSql, characterId)
	if err := row.Scan(&name); err != nil {
		return "", fmt.Errorf("row.Scan faild: %w", err)
	}
	return name, nil
}

func countPerRarity(rarity int) (int, error) {
	const selectSql = "SELECT COUNT(*) FROM characters WHERE rarity = $1"
	var count int
	row := db.QueryRow(selectSql, rarity)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("row.Scan faild: %w", err)
	}
	return count, nil
}

func selectCharacterList(token string) ([]Character, error) {
	var characters []Character
	const selectSql = `
SELECT UOC.id, C.id, C.name, UOC.level
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE U.digest_token = $1
`
	digestToken := hash(token)
	if _, err := selectUserId(token); err != nil {
		return nil, fmt.Errorf("selectUserId faild: %w", err)
	}

	rows, err := db.Query(selectSql, digestToken)
	if err != nil {
		return nil, fmt.Errorf("db.Query faild: %w", err)
	}
	for rows.Next() {
		var c Character
		if err := rows.Scan(&c.UserCharacterId, &c.CharacterId, &c.Name, &c.Level); err != nil {
			return nil, fmt.Errorf("rows.Scan faild: %w", err)
		}
		characters = append(characters, c)
	}
	return characters, nil
}

func selectCalorieByUserCharacterId(userCharacterId int) (int, error) {
	const selectSql = `
SELECT C.calorie
FROM user_ownership_characters AS UOC
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE UOC.id = $1
`
	row := db.QueryRow(selectSql, userCharacterId)
	var calorie int
	if err := row.Scan(&calorie); err != nil {
		return 0, fmt.Errorf("row.Scan faild: %w", err)
	}
	return calorie, nil
}

func selectExperience(userCharacterId int) (int, error) {
	const selectSql = `SELECT experience FROM user_ownership_characters WHERE id = $1`
	row := db.QueryRow(selectSql, userCharacterId)
	var experience int
	if err := row.Scan(&experience); err != nil {
		return 0, fmt.Errorf("row.Scan faild: %w", err)
	}
	return experience, nil
}

func calculateExperience(level int) int {
	return (level ^ 2) * 100
}

func calculateLevel(experience int) int {
	return int(math.Floor(math.Sqrt(float64(experience)) / 10.0))
}

func composeCharacter(baseUserCharacterId, materialUserCharacterId int) (*sql.Tx, int, error) {
	calorie, err := selectCalorieByUserCharacterId(materialUserCharacterId)
	if err != nil {
		return nil, 0, fmt.Errorf("selectCalorieByUserCharacterId faild: %w", err)
	}
	experience, err := selectExperience(baseUserCharacterId)
	if err != nil {
		return nil, 0, fmt.Errorf("selectExperience faild: %w", err)
	}
	newExperience := experience + calorie
	newLevel := calculateLevel(newExperience)

	tx, err := db.Begin()
	if err != nil {
		return nil, 0, fmt.Errorf("db.Begin faild: %w", err)
	}
	if err := updateCharacter(tx, baseUserCharacterId, newLevel, newExperience); err != nil {
		return tx, 0, fmt.Errorf("updateCharacter faild: %w", err)
	}
	if err := deleteCharacter(tx, materialUserCharacterId); err != nil {
		return tx, 0, fmt.Errorf("deleteCharacter faild: %w", err)
	}
	return tx, newLevel, nil
}

func updateCharacter(tx *sql.Tx, userCharacterId, level, experience int) error {
	const updateSql = `UPDATE user_ownership_characters SET level = $1, experience = $2 WHERE id = $3`
	if _, err := tx.Exec(updateSql, level, experience, userCharacterId); err != nil {
		return fmt.Errorf("tx.Exec faild: %w", err)
	}
	return nil
}

func deleteCharacter(tx *sql.Tx, userCharacterId int) error {
	const deleteSql = `DELETE FROM user_ownership_characters WHERE id = $1`
	if _, err := tx.Exec(deleteSql, userCharacterId); err != nil {
		return fmt.Errorf("tx.Exec faild: %w", err)
	}
	return nil
}

func createPutCharacterComposeResponse(userCharacterId, level int) (PutCharacterComposeResponse, error) {
	var jsonResponse PutCharacterComposeResponse
	const selectSql = `
SELECT UOC.id, C.id, C.name
FROM user_ownership_characters AS UOC
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE UOC.id = $1
`
	row := db.QueryRow(selectSql, userCharacterId)
	if err := row.Scan(&jsonResponse.UserCharacterId, &jsonResponse.CharacterId, &jsonResponse.Name); err != nil {
		return jsonResponse, fmt.Errorf("row.Scan faild: %w", err)
	}
	jsonResponse.Level = level
	return jsonResponse, nil
}
