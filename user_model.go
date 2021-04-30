package ca_game_api

import (
	"fmt"
)

type User struct {
	id          int
	name        string
	digestToken string
	createdAt   string
	updatedAt   string
}

func insertUser(user User) error {
	const createSql = `INSERT INTO users (name, digest_token) VALUES ($1, $2);`
	if _, err := db.Exec(createSql, user.name, user.digestToken); err != nil {
		return fmt.Errorf("db.Exec failed: %w", err)
	}
	return nil
}

func selectUserByToken(token string) (User, error) {
	const selectSql = `SELECT * FROM users WHERE digest_token = $1`
	digestToken := hash(token)

	var user User
	row := db.QueryRow(selectSql, digestToken)
	if err := row.Scan(&user.id, &user.name, &user.digestToken, &user.createdAt, &user.updatedAt); err != nil {
		return user, fmt.Errorf("row.Scan failed: %w", err)
	}
	return user, nil
}

func selectUserId(token string) (int, error) {
	const selectSql = `SELECT id FROM users WHERE digest_token = $1`
	digestToken := hash(token)
	row := db.QueryRow(selectSql, digestToken)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("row.Scan failed: %w", err)
	}
	return id, nil
}

func selectUserIdByUserCharacterId(userCharacterId int) (int, error) {
	const selectSql = `SELECT user_id FROM user_ownership_characters WHERE id = $1`
	row := db.QueryRow(selectSql, userCharacterId)
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("row.Scan failed: %w", err)
	}
	return id, nil
}

func updateUser(user User) error {
	const updateSql = `UPDATE users SET name = $1 WHERE digest_token = $2`
	if _, err := db.Exec(updateSql, user.name, user.digestToken); err != nil {
		return fmt.Errorf("db.Exec failed: %w", err)
	}
	return nil
}
