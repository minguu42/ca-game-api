package ca_game_api

import (
	"database/sql"
	"fmt"
)

func composeUserCharacter(tx *sql.Tx, baseUserCharacter, materialUserCharacter UserCharacter) error {
	baseUserCharacter.experience = baseUserCharacter.experience + materialUserCharacter.character.calorie

	if err := updateUserCharacter(tx, baseUserCharacter); err != nil {
		return fmt.Errorf("userCharacter.updateUser failed: %w", err)
	}

	if err := deleteUserCharacter(tx, materialUserCharacter.id); err != nil {
		return fmt.Errorf("materialUserCharacter.deleteUserCharacter failed: %w", err)
	}
	return nil
}
