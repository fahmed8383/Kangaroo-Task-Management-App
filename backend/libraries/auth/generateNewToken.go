package auth

import (
	"math/rand"
	"time"
)

// GenerateToken creates a random token that can be used for various forms of authorizations
func GenerateToken(size int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()")
	rand.Seed(time.Now().Unix())
	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
