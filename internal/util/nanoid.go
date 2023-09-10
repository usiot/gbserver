package util

import "crypto/rand"

const (
	// nanoidAlphabet is the default alphabet used by New().
	nanoidAlphabet = "-=0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultSize    = 21
)

var defaultNanoId = []byte("000000000000000000000")

func NanoId(l ...int) []byte {
	var size int = defaultSize
	if len(l) > 0 || l[0] > 0 {
		size = l[0]
	}

	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
		return defaultNanoId
	}

	nanoid := make([]byte, size)
	for i := 0; i < size; i++ {
		nanoid[i] = nanoidAlphabet[bytes[i]&63]
	}

	return nanoid
}
