package base62

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateCode(n uint8) (string, error) {
	if n == 0 {
		panic("base62: code length must be greater than zero")
	}

	code := make([]byte, n)
	max := big.NewInt(int64(len(alphabet)))

	for i := range code {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		code[i] = alphabet[index.Int64()]
	}

	return string(code), nil
}
