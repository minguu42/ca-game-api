package ca_game_api

import (
	"database/sql"
	"fmt"
	"time"
)

type gachaResult struct {
	id         int
	user       *User
	character  *Character
	experience int
	createdAt  time.Time
}

func insertResult(tx *sql.Tx, result gachaResult) error {
	const insertSql = "INSERT INTO gacha_results (user_id, character_id, experience) VALUES ($1, $2, $3)"
	if _, err := tx.Exec(insertSql, result.user.id, result.character.id, result.experience); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}
	return nil
}

func insertGachaResults(tx *sql.Tx, results []gachaResult) error {
	for _, result := range results {
		if err := insertResult(tx, result); err != nil {
			return fmt.Errorf("insertUser failed: %w", err)
		}

		if err := insertUserCharacter(tx, result.user.id, result.character.id, result.experience); err != nil {
			return fmt.Errorf("insertUserCharacter failed: %w", err)
		}
	}
	return nil
}
