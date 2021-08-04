package ca_game_api

import (
	"database/sql"
	"fmt"
)

func composeUserCharacter(tx *sql.Tx, baseUserCharacter, materialUserCharacter UserCharacter) (int, error) {
	baseUserCharacter.experience = baseUserCharacter.experience + materialUserCharacter.character.calorie

	if err := updateUserCharacter(tx, baseUserCharacter); err != nil {
		return 0, fmt.Errorf("userCharacter.updateUser failed: %w", err)
	}

	if err := deleteUserCharacter(tx, materialUserCharacter.id); err != nil {
		return 0, fmt.Errorf("materialUserCharacter.deleteUserCharacter failed: %w", err)
	}
	return baseUserCharacter.experience, nil
}
