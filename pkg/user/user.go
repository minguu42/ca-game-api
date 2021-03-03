package user

import (
	"database/sql"
	"github.com/minguu42/ca-game-api/pkg/helper"
	"log"
	"net/http"
)

func Insert(db *sql.DB, name string, token string) error {
	const createSql = "INSERT INTO users (name, digest_token) VALUES (?, ?)"
	digestToken := helper.HashToken(token)
	if _, err := db.Exec(createSql, name, digestToken); err != nil {
		return err
	}
	return nil
}

func GetName(db *sql.DB, token string) (string, error) {
	const selectSql = "SELECT name FROM users WHERE digest_token = ?"
	digestToken := helper.HashToken(token)
	var name string
	row := db.QueryRow(selectSql, digestToken)
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func GetId(db *sql.DB, token string) (int, error) {
	const selectSql = "SELECT id FROM users WHERE digest_token = ?"
	digestToken := helper.HashToken(token)
	row := db.QueryRow(selectSql, digestToken)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

type Character struct {
	UserCharacterId string `json:"userCharacterID"`
	CharacterId string `json:"characterID"`
	Name string `json:"name"`
}

func GetCharacterList(db *sql.DB, token string) ([]Character, error) {
	var characters []Character
	const selectSql = `
SELECT UOC.id, C.id, C.name
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE U.digest_token = ?
`
	digestToken := helper.HashToken(token)
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

func Update(db *sql.DB, token, newName string, w http.ResponseWriter) {
	const updateSql = "UPDATE users SET name = ? WHERE digest_token = ?"
	digestToken := helper.HashToken(token)
	if _, err := db.Exec(updateSql, newName, digestToken); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Update user error:", err)
	}
}
