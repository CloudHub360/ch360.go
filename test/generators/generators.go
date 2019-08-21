package generators

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

}

func Bytes() []byte {
	token := make([]byte, 100)
	rand.Read(token)
	return token
}

const lettersAndNumbers = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func String(root string) string {
	b := make([]byte, 8)

	for i := range b {
		b[i] = lettersAndNumbers[rand.Intn(len(lettersAndNumbers))]
	}
	return root + "-" + string(b)
}

func Bool() bool {
	if i := rand.Intn(100); i > 50 {
		return true
	} else {
		return false
	}
}

func Int() int {
	return rand.Int()
}
