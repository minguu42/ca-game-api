package ca_game_api

import (
	"fmt"
	"time"
)

type User struct {
	id          int
	name        string
	digestToken string
	createdAt   time.Time
	updatedAt   time.Time
}

func insertUser(name, digestToken string) error {
	const query = `INSERT INTO users (name, digest_token) VALUES ($1, $2);`
	if _, err := db.Exec(query, name, digestToken); err != nil {
		return fmt.Errorf("db.Exec failed: %v", err)
	}
	return nil
}

func getUserByDigestToken(digestToken string) (User, error) {
	const query = `SELECT id, name, digest_token, created_at, updated_at FROM users WHERE digest_token = $1`
	row := db.QueryRow(query, digestToken)

	var user User
	if err := row.Scan(&user.id, &user.name, &user.digestToken, &user.createdAt, &user.updatedAt); err != nil {
		return User{}, fmt.Errorf("row.Scan failed: %w", err)
	}
	return user, nil
}

func getUserById(id int) (User, error) {
	const query = `SELECT id, name, digest_token, created_at, updated_at FROM users WHERE id = $1`
	row := db.QueryRow(query, id)

	var user User
	if err := row.Scan(&user.id, &user.name, &user.digestToken, &user.createdAt, &user.updatedAt); err != nil {
		return User{}, fmt.Errorf("row.Scan failed: %w", err)
	}
	return user, nil
}

func updateUser(user User) error {
	const query = `UPDATE users SET name = $1 WHERE digest_token = $2`
	if _, err := db.Exec(query, user.name, user.digestToken); err != nil {
		return fmt.Errorf("db.Exec failed: %w", err)
	}
	return nil
}

func selectUserRanking() ([]UserRankingJson, error) {
	const query = `
SELECT U.name, SUM(UOC.experience * C.base_power) AS sumPower
FROM user_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
GROUP BY U.id
ORDER BY sumPower DESC
LIMIT 3
`
	var rankings []UserRankingJson
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("db.Query failed: %w", err)
	}
	for rows.Next() {
		var ranking UserRankingJson
		if err := rows.Scan(&ranking.Name, &ranking.SumPower); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		rankings = append(rankings, ranking)
	}
	return rankings, nil
}
