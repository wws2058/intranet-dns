package utils

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/google/uuid"
)

var (
	DefaultTimeFormat = fmt.Sprintf("%s %s", time.DateOnly, time.TimeOnly)
	alphanumeric      = []byte("0123456789abcdefghijklmnopqrstuvwxyz")
)

func GenUUID() string {
	return uuid.NewString()
}

// a random string of arbitrary length
func GenRandStr(n int) string {
	b := make([]byte, n)
	end := len(alphanumeric)
	for i := 0; i < n; i++ {
		b[i] = alphanumeric[rand.Intn(end)]
	}
	return string(b)
}

// subitem is in slice
func Contains[T comparable](as []T, sub T) bool {
	for _, v := range as {
		if v == sub {
			return true
		}
	}
	return false
}
