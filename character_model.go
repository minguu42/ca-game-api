package ca_game_api

import (
	"database/sql"
	"log"
	"math"
	"net/http"
)

type Character struct {
	UserCharacterId string `json:"userCharacterID"`
	CharacterId     string `json:"characterID"`
	Name            string `json:"name"`
	Level           int    `json:"level"`
}

func selectCharacterName(db *sql.DB, characterId int) (string, error) {
	const selectSql = "SELECT name FROM characters WHERE id = $1"
	var name string
	row := db.QueryRow(selectSql, characterId)
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func countPerRarity(db *sql.DB, rarity int, w http.ResponseWriter) (int, error) {
	const selectSql = "SELECT COUNT(*) FROM characters WHERE rarity = $1"
	var count int
	row := db.QueryRow(selectSql, rarity)
	if err := row.Scan(&count); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return 0, err
	}
	return count, nil
}

func selectCharacterList(token string, w http.ResponseWriter) ([]Character, error) {
	log.Println("INFO START selectCharacterList")
	var characters []Character
	const selectSql = `
SELECT UOC.id, C.id, C.name, UOC.level
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE U.digest_token = $1
`
	digestToken := HashToken(token)
	if _, err := selectUserId(token, w); err != nil {
		return nil, err
	}

	rows, err := db.Query(selectSql, digestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return nil, err
	}
	for rows.Next() {
		var c Character
		if err := rows.Scan(&c.UserCharacterId, &c.CharacterId, &c.Name, &c.Level); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("ERROR Return 500:", err)
			return nil, err
		}
		characters = append(characters, c)
	}
	log.Println("INFO END selectCharacterList")
	return characters, nil
}

func selectCalorieByUserCharacterId(db *sql.DB, userCharacterId int, w http.ResponseWriter) (int, error) {
	const selectSql = `
SELECT C.calorie
FROM user_ownership_characters AS UOC
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE UOC.id = $1
`
	row := db.QueryRow(selectSql, userCharacterId)
	var calorie int
	if err := row.Scan(&calorie); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return 0, err
	}
	return calorie, nil
}

func selectExperience(db *sql.DB, userCharacterId int, w http.ResponseWriter) (int, error) {
	const selectSql = `SELECT experience FROM user_ownership_characters WHERE id = $1`
	row := db.QueryRow(selectSql, userCharacterId)
	var experience int
	if err := row.Scan(&experience); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return 0, err
	}
	return experience, nil
}

func calculateExperience(level int) int {
	return (level ^ 2) * 100
}

func calculateLevel(experience int) int {
	return int(math.Floor(math.Sqrt(float64(experience)) / 10.0))
}

func composeCharacter(baseUserCharacterId, materialUserCharacterId int, w http.ResponseWriter) (*sql.Tx, int, error) {
	log.Println("INFO START composeCharacter")
	calorie, err := selectCalorieByUserCharacterId(db, materialUserCharacterId, w)
	if err != nil {
		return nil, 0, err
	}
	experience, err := selectExperience(db, baseUserCharacterId, w)
	if err != nil {
		return nil, 0, err
	}
	newExperience := experience + calorie
	newLevel := calculateLevel(newExperience)

	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return nil, 0, err
	}
	if err := updateCharacter(tx, baseUserCharacterId, newLevel, newExperience, w); err != nil {
		return tx, 0, err
	}
	if err := deleteCharacter(tx, materialUserCharacterId, w); err != nil {
		return tx, 0, err
	}
	log.Println("INFO END composeCharacter")
	return tx, newLevel, nil
}

func updateCharacter(tx *sql.Tx, userCharacterId, level, experience int, w http.ResponseWriter) error {
	const updateSql = `UPDATE user_ownership_characters SET level = $1, experience = $2 WHERE id = $3`
	if _, err := tx.Exec(updateSql, level, experience, userCharacterId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return err
	}
	return nil
}

func deleteCharacter(tx *sql.Tx, userCharacterId int, w http.ResponseWriter) error {
	const deleteSql = `DELETE FROM user_ownership_characters WHERE id = $1`
	if _, err := tx.Exec(deleteSql, userCharacterId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return err
	}
	return nil
}

func createPutCharacterComposeResponse(userCharacterId, level int, w http.ResponseWriter) (PutCharacterComposeResponse, error) {
	var jsonResponse PutCharacterComposeResponse
	const selectSql = `
SELECT UOC.id, C.id, C.name
FROM user_ownership_characters AS UOC
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE UOC.id = $1
`
	row := db.QueryRow(selectSql, userCharacterId)
	if err := row.Scan(&jsonResponse.UserCharacterId, &jsonResponse.CharacterId, &jsonResponse.Name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return jsonResponse, err
	}
	jsonResponse.Level = level
	return jsonResponse, nil
}
