package data

import (
	"math/rand"
)

func GenerateAccessCode() string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := ""

	for i := 0; i < 5; i++ {
		idx := rand.Intn(len(characters))
		code += string(characters[idx])
	}
	return code
}

