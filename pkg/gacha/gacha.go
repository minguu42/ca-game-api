package gacha

import (
	"database/sql"
	"github.com/minguu42/ca-game-api/pkg/character"
	"github.com/minguu42/ca-game-api/pkg/helper"
	"strconv"
)

type Result struct {
	CharacterId string `json:"characterID"`
	Name        string `json:"name"`
}

func ApplyGachaResult(db *sql.DB, userId, CharacterId int) error {
	const insertSql = "INSERT INTO gacha_results (user_id, character_id) VALUES (?, ?)"
	if _, err := db.Exec(insertSql, userId, CharacterId); err != nil {
		return err
	}
	const createSql = "INSERT INTO user_ownership_characters (user_id, character_id) VALUES (?, ?)"
	if _, err := db.Exec(createSql, userId, CharacterId); err != nil {
		return err
	}
	return nil
}

func Draw(db *sql.DB, userId, times int) ([]Result, error) {
	var results []Result
	for i := 0; i < times; i++ {
		rarity3CharacterNum, err := character.CountPerRarity(db, 3)
		if err != nil {
			return nil, err
		}
		rarity4CharacterNum, err := character.CountPerRarity(db, 4)
		if err != nil {
			return nil, err
		}
		rarity5CharacterNum, err := character.CountPerRarity(db, 5)
		if err != nil {
			return nil, err
		}
		var characterId int
		switch selectedRarity := helper.SelectRarity(); selectedRarity{
		case 3:
			characterId = helper.SelectCharacterId(rarity3CharacterNum) +  + 30000000
		case 4:
			characterId = helper.SelectCharacterId(rarity4CharacterNum) + 40000000
		case 5:
			characterId = helper.SelectCharacterId(rarity5CharacterNum) + 50000000
		}
		name, err := character.GetName(db, characterId)
		if err != nil {
			return nil, err
		}
		if err := ApplyGachaResult(db, userId, characterId); err != nil {
			return nil, err
		}
		results = append(results, Result{
			CharacterId: strconv.Itoa(characterId),
			Name:        name,
		})
	}
	return results, nil
}