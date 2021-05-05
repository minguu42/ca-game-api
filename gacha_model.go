package ca_game_api

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type gachaResult struct {
	id        int
	user      *User
	character *Character
	level     int
	createdAt time.Time
}

func (gachaResult gachaResult) insert(tx *sql.Tx) error {
	const insertSql = "INSERT INTO gacha_results (user_id, character_id, level) VALUES ($1, $2, $3)"
	if _, err := tx.Exec(insertSql, gachaResult.user.id, gachaResult.character.id, gachaResult.level); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}
	return nil
}

type UserOwnCharacter struct {
	id         int
	user       *User
	character  *Character
	level      int
	experience int
	createdAt  time.Time
	updatedAt  time.Time
}

func (userOwnCharacter UserOwnCharacter) insert(tx *sql.Tx) error {
	const stmt = `INSERT INTO user_ownership_characters (user_id, character_id, level, experience) VALUES ($1, $2, $3, $4)`
	if _, err := tx.Exec(stmt, userOwnCharacter.user.id, userOwnCharacter.character.id, userOwnCharacter.level, userOwnCharacter.level); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}
	return nil
}

func decideRarity() int {
	if num := rand.Intn(1000) + 1; num >= 900 {
		return 5
	} else if num >= 600 {
		return 4
	} else {
		return 3
	}
}

func decideCharacterId() (int, error) {
	rarity3Num, err := countCharactersByRarity(3)
	if err != nil {
		return 0, fmt.Errorf("countCharactersByRarity failed: %w", err)
	}
	rarity4Num, err := countCharactersByRarity(4)
	if err != nil {
		return 0, fmt.Errorf("countCharactersByRarity failed: %w", err)
	}
	rarity5Num, err := countCharactersByRarity(5)
	if err != nil {
		return 0, fmt.Errorf("countCharactersByRarity failed: %w", err)
	}

	var characterId int
	switch rarity := decideRarity(); rarity {
	case 3:
		characterId = rand.Intn(rarity3Num) + 30000001
	case 4:
		characterId = rand.Intn(rarity4Num) + 40000001
	case 5:
		characterId = rand.Intn(rarity5Num) + 50000001
	}
	return characterId, nil
}

func decideCharacterLevel() int {
	return rand.Intn(10) + 1
}

func decideGachaResults(user User, times int) ([]gachaResult, error) {
	results := make([]gachaResult, 0, times)

	for i := 0; i < times; i++ {
		rand.Seed(time.Now().UnixNano())

		characterId, err := decideCharacterId()
		if err != nil {
			return nil, fmt.Errorf("decideCharacterId failed: %w", err)
		}
		character, err := getCharacterById(characterId)
		if err != nil {
			return nil, fmt.Errorf("getCharacterById failed: %w", err)
		}
		level := decideCharacterLevel()

		results = append(results, gachaResult{
			user:      &user,
			character: &character,
			level:     level,
		})
	}
	return results, nil
}

func storeGachaResults(tx *sql.Tx, results []gachaResult) error {
	for _, result := range results {
		userOwnCharacter := UserOwnCharacter{
			user:       result.user,
			character:  result.character,
			level:      result.level,
			experience: calculateExperience(result.level),
		}
		if err := result.insert(tx); err != nil {
			return fmt.Errorf("result.insert failed: %w", err)
		}
		if err := userOwnCharacter.insert(tx); err != nil {
			return fmt.Errorf("userOwnCharacter.insert failed: %w", err)
		}
	}
	return nil
}
