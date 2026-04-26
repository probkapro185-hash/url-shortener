package storage

import (
	"math/rand"
	"time"
)

const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode(lenght int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, lenght)
	for i := range result {
		result[i] = char[rand.Intn(len(char))]
	}
	return string(result)
}
