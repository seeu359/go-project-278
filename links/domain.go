package links

import (
	"net/url"
	"os"
)

func GetShortURL(shortName string) (string, error) {
	baseName := os.Getenv("BASE_NAME")
	res, err := url.JoinPath(baseName, shortName)
	if err != nil {
		return "", err
	}
	return res, nil
}
