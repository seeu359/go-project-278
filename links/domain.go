package links

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"os"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetShortURL(shortName string) (string, error) {
	baseName := os.Getenv("BASE_NAME")
	res, err := url.JoinPath(baseName, shortName)
	if err != nil {
		return "", err
	}
	return res, nil
}

func RandomShortName(length int) (string, error) {
	result := make([]byte, length)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}

		result[i] = letters[n.Int64()]
	}

	return string(result), nil
}
