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

func draw(tx *sql.Tx, userId int, times int) ([]ResultJson, error) {
	var results []ResultJson

	for i := 0; i < times; i++ {
		rand.Seed(time.Now().UnixNano())
		characterId, err := decideCharacterId()
		if err != nil {
			return nil, fmt.Errorf("decideCharacterId failed: %w", err)
		}
		characterLevel := decideCharacterLevel()
		characterExperience := calculateExperience(characterLevel)
		character, err := getCharacterById(characterId)
		if err != nil {
			return nil, fmt.Errorf("getCharacterById failed: %w", err)
		}

		results = append(results, ResultJson{
			CharacterId: characterId,
			Name:        character.name,
		})

		if err := insertResult(tx, userId, characterId, characterLevel, characterExperience); err != nil {
			return nil, fmt.Errorf("insertResult fiald: %w", err)
		}
	}
	return results, nil
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

func insertResult(tx *sql.Tx, userId, characterId, characterLevel, characterExperience int) error {
	const insertSql = "INSERT INTO gacha_results (user_id, character_id, level) VALUES ($1, $2, $3)"
	if _, err := tx.Exec(insertSql, userId, characterId, characterLevel); err != nil {
		return fmt.Errorf("tx.Exec faild: %w", err)
	}
	const createSql = "INSERT INTO user_ownership_characters (user_id, character_id, level, experience) VALUES ($1, $2, $3, $4)"
	if _, err := tx.Exec(createSql, userId, characterId, characterLevel, characterExperience); err != nil {
		return fmt.Errorf("tx.Exec faild: %w", err)
	}
	return nil
}
