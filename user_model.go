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

type UserRankingInfo struct {
	Id       string `json:"userID"`
	Name     string `json:"name"`
	SumPower string `json:"sumPower"`
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

func selectUserRanking() ([]UserRankingInfo, error) {
	var users []UserRankingInfo
	const selectSql = `
SELECT U.id, U.name, SUM(UOC.level * C.power) AS sumPower
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
GROUP BY U.id
ORDER BY sumPower DESC
LIMIT 3
`
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("db.Query faild: %w", err)
	}
	for rows.Next() {
		var user UserRankingInfo
		if err := rows.Scan(&user.Id, &user.Name, &user.SumPower); err != nil {
			return nil, fmt.Errorf("rows.Scan faild: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
