package ca_game_api

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashToken(token string) string {
	digestTokenByte := sha256.Sum256([]byte(token))
	digestToken := hex.EncodeToString(digestTokenByte[:])
	return digestToken
}