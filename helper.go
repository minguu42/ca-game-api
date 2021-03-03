package ca_game_api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}

func HashToken(token string) string {
	digestTokenByte := sha256.Sum256([]byte(token))
	digestToken := hex.EncodeToString(digestTokenByte[:])
	return digestToken
}
