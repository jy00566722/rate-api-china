package utils

import (
	"math/rand"
	"time"
)

func GenerateVerifyCode(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	digits := "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[r.Intn(len(digits))]
	}
	return string(code)
}
