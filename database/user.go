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
	row := db.QueryRow(selectSql, digestToken)
	var name string
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