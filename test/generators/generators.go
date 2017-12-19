package generators

import (
	"math/rand"
	"time"
)

func Bytes() []byte {
	token := make([]byte, 100)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(token)
	return token
}

const lettersAndNumbers = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func String(root string) string {
	b := make([]byte, 8)
	rand.Seed(time.Now().UTC().UnixNano())

	for i := range b {
		b[i] = lettersAndNumbers[rand.Intn(len(lettersAndNumbers))]
	}
	return root + "-" + string(b)
}
