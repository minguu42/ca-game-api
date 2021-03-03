package ca_game_api

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"
)

type Result struct {
	CharacterId string `json:"characterID"`
	Name        string `json:"name"`
}

func Draw(db *sql.DB, userId, times int) ([]Result, error) {
	var results []Result

	rarity3SumNum, err := countPerRarity(db, 3)
	if err != nil {
		return nil, err
	}
	rarity4SumNum, err := countPerRarity(db, 4)
	if err != nil {
		return nil, err
	}
	rarity5SumNum, err := countPerRarity(db, 5)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	for i := 0; i < times; i++ {
		characterId := selectCharacterId(rarity3SumNum, rarity4SumNum, rarity5SumNum)
		name, err := selectCharacterName(db, characterId)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
		results = append(results, Result{
			CharacterId: strconv.Itoa(characterId),
			Name:        name,
		})
		if err := applyResult(tx, userId, characterId); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return results, nil
}

func selectRarity() int {
	rand.Seed(time.Now().UnixNano())
	if num := rand.Intn(1000) + 1; num >= 900 {
		return 5
	} else if num >= 600 {
		return 4
	} else {
		return 3
	}
}

func selectCharacterId(rarity3SumNum, rarity4SumNum, rarity5SumNum int) int {
	var characterId int
	switch rarity := selectRarity(); rarity {
	case 3:
		characterId = rand.Intn(rarity3SumNum) + 30000001
	case 4:
		characterId = rand.Intn(rarity4SumNum) + 40000001
	case 5:
		characterId = rand.Intn(rarity5SumNum) + 50000001
	}
	return characterId
}

func applyResult(tx *sql.Tx, userId, CharacterId int) error {
	const insertSql = "INSERT INTO gacha_results (user_id, character_id) VALUES (?, ?)"
	if _, err := tx.Exec(insertSql, userId, CharacterId); err != nil {
		return err
	}
	const createSql = "INSERT INTO user_ownership_characters (user_id, character_id) VALUES (?, ?)"
	if _, err := tx.Exec(createSql, userId, CharacterId); err != nil {
		return err
	}
	return nil
}
