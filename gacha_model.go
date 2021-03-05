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

func draw(db *sql.DB, xToken string, times int, w http.ResponseWriter) ([]Result, error) {
	log.Println("INFO START draw")
	var results []Result

	userId, err := selectUserId(db, xToken, w)
	if err != nil {
		return nil, err
	}

	rarity3SumNum, err := countPerRarity(db, 3, w)
	if err != nil {
		return nil, err
	}
	rarity4SumNum, err := countPerRarity(db, 4, w)
	if err != nil {
		return nil, err
	}
	rarity5SumNum, err := countPerRarity(db, 5, w)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return nil, err
	}
	for i := 0; i < times; i++ {
		log.Printf("INFO START Draw once gach (%v / %v) \n", i+1, times)
		characterId := decideOutputCharacterId(rarity3SumNum, rarity4SumNum, rarity5SumNum)
		name, err := selectCharacterName(db, characterId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("ERROR Return 500:", err)
			if err := tx.Rollback(); err != nil {
				log.Println("ERROR Rollback error:", err)
				return nil, err
			}
			return nil, err
		}
		results = append(results, Result{
			CharacterId: strconv.Itoa(characterId),
			Name:        name,
		})
		if err := applyResult(tx, userId, characterId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("ERROR Return 500:", err)
			if err := tx.Rollback(); err != nil {
				log.Println("ERROR Rollback error:", err)
				return nil, err
			}
			return nil, err
		}
		log.Printf("INFO END Draw once gach (%v / %v) \n", i+1, times)
	}
	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return nil, err
	}
	log.Println("INFO Settle gacha result")
	log.Println("INFO END draw")
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

func decideOutputCharacterId(rarity3SumNum, rarity4SumNum, rarity5SumNum int) int {
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
