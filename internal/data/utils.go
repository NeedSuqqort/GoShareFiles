package data

import (
	"math/rand"
)

func GenerateAccessCode() string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := ""
	for range 5 {
		idx := rand.Intn(len(characters))
		code += string(characters[idx])
	}
	return code
}
