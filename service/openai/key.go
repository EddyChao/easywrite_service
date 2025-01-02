package openai

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateKey(length int) string {
	rand.Seed(time.Now().UnixNano())

	key := make([]byte, length)
	for i := 0; i < length; i++ {
		key[i] = charset[rand.Intn(len(charset))]
	}

	return "sk-" + string(key)
}
