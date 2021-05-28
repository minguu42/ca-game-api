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

func (gachaResult gachaResult) insert(tx *sql.Tx) error {
	const insertSql = "INSERT INTO gacha_results (user_id, character_id, experience) VALUES ($1, $2, $3)"
	if _, err := tx.Exec(insertSql, gachaResult.user.id, gachaResult.character.id, gachaResult.experience); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}
	return nil
}

func storeGachaResults(tx *sql.Tx, results []gachaResult) error {
	for _, result := range results {
		userOwnCharacter := UserCharacter{
			user:       result.user,
			character:  result.character,
			experience: result.experience,
		}
		if err := result.insert(tx); err != nil {
			return fmt.Errorf("result.insertUser failed: %w", err)
		}
		if err := userOwnCharacter.insert(tx); err != nil {
			return fmt.Errorf("userOwnCharacter.insertUser failed: %w", err)
		}
	}
	return nil
}
