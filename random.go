package hermes

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var charsetLen = big.NewInt(int64(len(charset)))

func RandomString(n int) (string, error) {
	b := make([]byte, n)

	for i := range b {
		index, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}

		b[i] = charset[index.Int64()]
	}

	return string(b), nil
}
