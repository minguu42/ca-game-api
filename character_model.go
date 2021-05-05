package ca_game_api

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

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
		return character, fmt.Errorf("row.Scan failed: %w", err)
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

type UserCharacter struct {
	id         int
	user       *User
	character  *Character
	level      int
	experience int
	createdAt  time.Time
	updatedAt  time.Time
}

func getUserCharacterById(id int) (UserCharacter, error) {
	const query = `SELECT id, user_id, character_id, level, experience, created_at, updated_at FROM user_ownership_characters WHERE id = $1`
	row := db.QueryRow(query, id)
	var userCharacter UserCharacter
	var userId int
	var characterId int
	if err := row.Scan(&userCharacter.id, &userId, &characterId, &userCharacter.level, &userCharacter.experience, &userCharacter.createdAt, &userCharacter.updatedAt); err != nil {
		return userCharacter, fmt.Errorf("row.Scan failed: %w", err)
	}

	user, err := getUserById(userId)
	if err != nil {
		return userCharacter, fmt.Errorf("getUserById failed: %w", err)
	}
	character, err := getCharacterById(characterId)
	if err != nil {
		return userCharacter, fmt.Errorf("getCharacterById failed: %w", err)
	}
	userCharacter.user = &user
	userCharacter.character = &character
	return userCharacter, nil
}

func (userCharacter UserCharacter) insert(tx *sql.Tx) error {
	const query = `INSERT INTO user_ownership_characters (user_id, character_id, level, experience) VALUES ($1, $2, $3, $4)`
	if _, err := tx.Exec(query, userCharacter.user.id, userCharacter.character.id, userCharacter.level, userCharacter.level); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}
	return nil
}

func getUserOwnCharactersByToken(token string) ([]UserCharacter, error) {
	user, err := getUserByToken(token)
	if err != nil {
		return nil, fmt.Errorf("getUserByToken failed: %w", err)
	}

	const query = `
SELECT UOC.id, C.id, UOC.level, UOC.experience
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE U.digest_token = $1
`
	rows, err := db.Query(query, user.digestToken)
	if err != nil {
		return nil, fmt.Errorf("db.Query failed: %w", err)
	}
	var userCharacters []UserCharacter
	for rows.Next() {
		var userOwnCharacter UserCharacter
		var characterId int
		if err := rows.Scan(&userOwnCharacter.id, &characterId, &userOwnCharacter.level, &userOwnCharacter.experience); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		character, err := getCharacterById(characterId)
		if err != nil {
			return nil, fmt.Errorf("getCharacterById failed: %w", err)
		}

		userOwnCharacter.user = &user
		userOwnCharacter.character = &character
		userCharacters = append(userCharacters, userOwnCharacter)
	}
	return userCharacters, nil
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

func calculatePower(userCharacter UserCharacter) int {
	return userCharacter.level * userCharacter.character.basePower
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
