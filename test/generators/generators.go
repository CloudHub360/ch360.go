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
