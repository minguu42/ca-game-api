package database

import (
	"database/sql"
)

func InsertUser(db *sql.DB, name string, digestToken string) error {
	const createSql = "INSERT INTO users (name, digest_token) VALUES (?, ?)"
	if _, err := db.Exec(createSql, name, digestToken); err != nil {
		return err
	}
	return nil
}

func GetUserName(db *sql.DB, digestToken string) (string, error) {
	const selectSql = "SELECT name FROM users WHERE digest_token = ?"
	var name string
	row := db.QueryRow(selectSql, digestToken)
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

func GetUserId(db *sql.DB, digestToken string) (int, error) {
	const selectSql = "SELECT id FROM users WHERE digest_token = ?"
	row := db.QueryRow(selectSql, digestToken)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateUser(db *sql.DB, id int, newName string) error {
	const updateSql = "UPDATE users SET name = ? WHERE id = ?"
	if _, err := db.Exec(updateSql, newName, id); err != nil {
		return err
	}
	return nil
}

type Character struct {
	UserCharacterId string `json:"userCharacterID"`
	CharacterId string `json:"characterID"`
	Name string `json:"name"`
}

func GetCharacterList(db *sql.DB, digestToken string) ([]Character, error) {
	var characters []Character
	const selectSql = `
SELECT UOC.id, C.id, C.name
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE U.digest_token = ?
`
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