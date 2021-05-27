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
