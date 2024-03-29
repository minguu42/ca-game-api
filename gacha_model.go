package ca_game_api

import (
	"fmt"
	"math/rand"
	"time"
)

func decideRarity(probabilityOf5, probabilityOf4, probabilityOf3 int) int {
	if num := rand.Intn(probabilityOf5+probabilityOf4+probabilityOf3) + 1; num <= probabilityOf5 {
		return 5
	} else if num <= probabilityOf4 {
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
	switch rarity := decideRarity(10, 30, 60); rarity {
	case 3:
		characterId = rand.Intn(rarity3Num) + 30000001
	case 4:
		characterId = rand.Intn(rarity4Num) + 40000001
	case 5:
		characterId = rand.Intn(rarity5Num) + 50000001
	}
	return characterId, nil
}

func calculateCharacterExperience(level int) int {
	return level * level * 100
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
		experience := calculateCharacterExperience(rand.Intn(10) + 1)

		results = append(results, gachaResult{
			user:       &user,
			character:  &character,
			experience: experience,
		})
	}
	return results, nil
}
