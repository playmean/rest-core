package restcore

import "math/rand"

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandomString(length int) string {
	b := make([]rune, length)

	charsLen := len(chars)

	for i := range b {
		b[i] = chars[rand.Intn(charsLen)]
	}

	return string(b)
}
