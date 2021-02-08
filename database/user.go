package database

import (
	"database/sql"
)

func InsertUser(db *sql.DB, name string, digestToken string) error {
	const createSql = "INSERT INTO users (name, token_digest) VALUES (?, ?)"
	_, err := db.Exec(createSql, name, digestToken)
	if err != nil {
		return err
	}
	return nil
}

func GetUserName(db *sql.DB, digestToken string) (string, error) {
	const selectSql = "SELECT name FROM users WHERE token_digest = ?"
	row := db.QueryRow(selectSql, digestToken)
	var name string
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}