package ca_game_api

import (
	"database/sql"
	"log"
	"net/http"
)

func insertUser(db *sql.DB, name string, token string) error {
	const createSql = "INSERT INTO users (name, digest_token) VALUES (?, ?)"
	digestToken := HashToken(token)
	if _, err := db.Exec(createSql, name, digestToken); err != nil {
		return err
	}
	return nil
}

func selectUserName(db *sql.DB, token string) (string, error) {
	const selectSql = "SELECT name FROM users WHERE digest_token = ?"
	digestToken := HashToken(token)
	var name string
	row := db.QueryRow(selectSql, digestToken)
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func selectUserId(db *sql.DB, token string) (int, error) {
	const selectSql = "SELECT id FROM users WHERE digest_token = ?"
	digestToken := HashToken(token)
	row := db.QueryRow(selectSql, digestToken)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

type Character struct {
	UserCharacterId string `json:"userCharacterID"`
	CharacterId     string `json:"characterID"`
	Name            string `json:"name"`
}

func selectCharacterList(db *sql.DB, token string) ([]Character, error) {
	var characters []Character
	const selectSql = `
SELECT UOC.id, C.id, C.name
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE U.digest_token = ?
`
	digestToken := HashToken(token)
	rows, err := db.Query(selectSql, digestToken)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var c Character
		if err := rows.Scan(&c.UserCharacterId, &c.CharacterId, &c.Name); err != nil {
			return nil, err
		}
		characters = append(characters, c)
	}
	return characters, nil
}

func updateUser(db *sql.DB, token, newName string, w http.ResponseWriter) {
	const updateSql = "UPDATE users SET name = ? WHERE digest_token = ?"
	digestToken := HashToken(token)
	if _, err := db.Exec(updateSql, newName, digestToken); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR updateUser user error:", err)
	}
}
