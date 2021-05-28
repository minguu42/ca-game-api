package ca_game_api

import (
	"database/sql"
	"fmt"
	"time"
)

type UserCharacter struct {
	id         int
	user       *User
	character  *Character
	experience int
	createdAt  time.Time
	updatedAt  time.Time
}

func (userCharacter *UserCharacter) insert(tx *sql.Tx) error {
	const query = `INSERT INTO user_characters (user_id, character_id, experience) VALUES ($1, $2, $3)`
	if _, err := tx.Exec(query, userCharacter.user.id, userCharacter.character.id, userCharacter.experience); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}
	return nil
}

func getUserCharacterById(id int) (UserCharacter, error) {
	const query = `SELECT id, user_id, character_id, experience, created_at, updated_at FROM user_characters WHERE id = $1`
	row := db.QueryRow(query, id)
	var userCharacter UserCharacter
	var userId int
	var characterId int
	if err := row.Scan(&userCharacter.id, &userId, &characterId, &userCharacter.experience, &userCharacter.createdAt, &userCharacter.updatedAt); err != nil {
		return userCharacter, fmt.Errorf("row.Scan failed: %w", err)
	}

	user, err := getUserById(userId)
	if err != nil {
		return userCharacter, fmt.Errorf("getUserById failed: %w", err)
	}
	character, err := getCharacterById(characterId)
	if err != nil {
		return userCharacter, fmt.Errorf("getCharacterById failed: %w", err)
	}
	userCharacter.user = &user
	userCharacter.character = &character
	return userCharacter, nil
}

func getUserCharactersByToken(token string) ([]UserCharacter, error) {
	user, err := getUserByToken(token)
	if err != nil {
		return nil, fmt.Errorf("getUserByToken failed: %w", err)
	}

	const query = `
SELECT UOC.id, C.id, UOC.experience
FROM user_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
WHERE U.digest_token = $1
`
	rows, err := db.Query(query, user.digestToken)
	if err != nil {
		return nil, fmt.Errorf("db.Query failed: %w", err)
	}
	var userCharacters []UserCharacter
	for rows.Next() {
		var userOwnCharacter UserCharacter
		var characterId int
		if err := rows.Scan(&userOwnCharacter.id, &characterId, &userOwnCharacter.experience); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		character, err := getCharacterById(characterId)
		if err != nil {
			return nil, fmt.Errorf("getCharacterById failed: %w", err)
		}

		userOwnCharacter.user = &user
		userOwnCharacter.character = &character
		userCharacters = append(userCharacters, userOwnCharacter)
	}
	return userCharacters, nil
}

func (userCharacter *UserCharacter) update(tx *sql.Tx) error {
	const query = `UPDATE user_characters SET experience = $2 WHERE id = $1`
	if _, err := tx.Exec(query, userCharacter.id, userCharacter.experience); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}
	return nil
}

func (userCharacter *UserCharacter) delete(tx *sql.Tx) error {
	const query = `DELETE FROM user_characters WHERE id = $1`
	if _, err := tx.Exec(query, userCharacter.id); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}
	return nil
}

func (userCharacter *UserCharacter) compose(tx *sql.Tx, materialUserCharacter UserCharacter) error {
	userCharacter.experience = userCharacter.experience + materialUserCharacter.character.calorie
	if err := userCharacter.update(tx); err != nil {
		return fmt.Errorf("userCharacter.update failed: %w", err)
	}
	if err := materialUserCharacter.delete(tx); err != nil {
		return fmt.Errorf("materialUserCharacter.delete failed: %w", err)
	}
	return nil
}
