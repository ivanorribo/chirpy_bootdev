package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "" // might handle the error in the future
	}
	encodedStr := hex.EncodeToString(randomBytes)
	return encodedStr
}
