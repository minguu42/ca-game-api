package ca_game_api

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Result struct {
	CharacterId string `json:"characterID"`
	Name        string `json:"name"`
}

func draw(db *sql.DB, xToken string, times int, w http.ResponseWriter) ([]Result, error, *sql.Tx) {
	log.Println("INFO START draw")
	var results []Result

	userId, err := selectUserId(db, xToken, w)
	if err != nil {
		return nil, err, nil
	}

	rarity3SumNum, err := countPerRarity(db, 3, w)
	if err != nil {
		return nil, err, nil
	}
	rarity4SumNum, err := countPerRarity(db, 4, w)
	if err != nil {
		return nil, err, nil
	}
	rarity5SumNum, err := countPerRarity(db, 5, w)
	if err != nil {
		return nil, err, nil
	}

	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return nil, err, nil
	}
	for i := 0; i < times; i++ {
		log.Printf("INFO START Draw once gach (%v / %v) \n", i+1, times)
		rand.Seed(time.Now().UnixNano())
		characterId := decideOutputCharacterId(rarity3SumNum, rarity4SumNum, rarity5SumNum)
		characterLevel := decideCharacterLevel()
		characterExperience := calculateExperience(characterLevel)
		name, err := selectCharacterName(db, characterId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("ERROR Return 500:", err)
			return nil, err, tx
		}
		results = append(results, Result{
			CharacterId: strconv.Itoa(characterId),
			Name:        name,
		})

		if err := insertResult(tx, userId, characterId, characterLevel, characterExperience); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("ERROR Return 500:", err)
			return nil, err, tx
		}
		log.Printf("INFO END Draw once gach (%v / %v) \n", i+1, times)
	}
	log.Println("INFO END draw")
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
	const insertSql = "INSERT INTO gacha_results (user_id, character_id, level) VALUES (?, ?, ?)"
	if _, err := tx.Exec(insertSql, userId, characterId, characterLevel); err != nil {
		return err
	}
	const createSql = "INSERT INTO user_ownership_characters (user_id, character_id, level, experience) VALUES (?, ?, ?, ?)"
	if _, err := tx.Exec(createSql, userId, characterId, characterLevel, characterExperience); err != nil {
		return err
	}
	return nil
}
