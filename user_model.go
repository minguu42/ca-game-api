package ca_game_api

import (
	"fmt"
)

func insertUser(name string, token string) error {
	const createSql = `INSERT INTO users (name, digest_token) VALUES ($1, $2);`
	digestToken := hash(token)
	if _, err := db.Exec(createSql, name, digestToken); err != nil {
		return fmt.Errorf("db.Exec faild: %w", err)
	}
	return nil
}

func selectUserByToken(token string) (string, error) {
	const selectSql = `SELECT name FROM users WHERE digest_token = $1`
	digestToken := hash(token)

	var name string
	row := db.QueryRow(selectSql, digestToken)
	if err := row.Scan(&name); err != nil {
		return "", fmt.Errorf("row.Scan faild: %w", err)
	}
	return name, nil
}

func selectUserId(token string) (int, error) {
	const selectSql = `SELECT id FROM users WHERE digest_token = $1`
	digestToken := hash(token)
	row := db.QueryRow(selectSql, digestToken)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("row.Scan faild: %w", err)
	}
	return id, nil
}

func selectUserIdByUserCharacterId(userCharacterId int) (int, error) {
	const selectSql = `SELECT user_id FROM user_ownership_characters WHERE id = $1`
	row := db.QueryRow(selectSql, userCharacterId)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("row.Scan faild: %w", err)
	}
	return id, nil
}

func updateUser(token, newName string) error {
	const updateSql = `UPDATE users SET name = $1 WHERE digest_token = $2`
	digestToken := hash(token)
	if _, err := db.Exec(updateSql, newName, digestToken); err != nil {
		return fmt.Errorf("db.Exec faild: %w", err)
	}
	return nil
}
