package gen

import (
	"math/rand"
	"sync"
	"time"
)

const (
	letters = "qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890-_"
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu         = &sync.Mutex{}
)

func LatinStr(length int) string {
	b := make([]byte, length)
	mu.Lock()
	defer mu.Unlock()
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(b)
}

func Intn(n int) int {
	mu.Lock()
	defer mu.Unlock()
	return seededRand.Intn(n)
}
