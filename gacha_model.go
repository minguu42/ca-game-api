package ca_game_api

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type gachaResult struct {
	id        int
	user      *User
	character *Character
	level     int
	createdAt time.Time
}

type Result struct {
	CharacterId string `json:"characterID"`
	Name        string `json:"name"`
}

func draw(xToken string, times int) ([]Result, error, *sql.Tx) {
	var results []Result

	userId, err := selectUserId(xToken)
	if err != nil {
		return nil, fmt.Errorf("selectUserid faild: %w", err), nil
	}

	rarity3SumNum, err := countPerRarity(3)
	if err != nil {
		return nil, fmt.Errorf("countPerRarity faild: %w", err), nil
	}
	rarity4SumNum, err := countPerRarity(4)
	if err != nil {
		return nil, fmt.Errorf("countPerRarity faild: %w", err), nil
	}
	rarity5SumNum, err := countPerRarity(5)
	if err != nil {
		return nil, fmt.Errorf("countPerRarity faild: %w", err), nil
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("db.Begin faild: %w", err), nil
	}
	for i := 0; i < times; i++ {
		rand.Seed(time.Now().UnixNano())
		characterId := decideOutputCharacterId(rarity3SumNum, rarity4SumNum, rarity5SumNum)
		characterLevel := decideCharacterLevel()
		characterExperience := calculateExperience(characterLevel)
		name, err := selectCharacterName(characterId)
		if err != nil {
			return nil, fmt.Errorf("selectCharacterName faild: %w", err), tx
		}
		results = append(results, Result{
			CharacterId: strconv.Itoa(characterId),
			Name:        name,
		})

		if err := insertResult(tx, userId, characterId, characterLevel, characterExperience); err != nil {
			return nil, fmt.Errorf("insertResult fiald: %w", err), tx
		}
	}
	return results, nil, tx
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

func decideOutputCharacterId(rarity3SumNum, rarity4SumNum, rarity5SumNum int) int {
	var characterId int
	switch rarity := decideRarity(); rarity {
	case 3:
		characterId = rand.Intn(rarity3SumNum) + 30000001
	case 4:
		characterId = rand.Intn(rarity4SumNum) + 40000001
	case 5:
		characterId = rand.Intn(rarity5SumNum) + 50000001
	}
	return characterId
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
