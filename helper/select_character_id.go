package helper

import (
	"math/rand"
	"time"
)

func SelectCharacterId(characterSum int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(characterSum) + 1
}
