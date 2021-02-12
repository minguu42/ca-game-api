package helper

import (
	"math/rand"
	"time"
)

func SelectRarity() int {
	rand.Seed(time.Now().UnixNano())
	if num := rand.Intn(1000) + 1; num >= 900 {
		return 5
	} else if num >= 600 {
		return 4
	} else {
		return 3
	}
}