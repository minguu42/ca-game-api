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