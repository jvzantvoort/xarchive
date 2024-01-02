package target

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func GetSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	bs := hash.Sum(nil)
	return fmt.Sprintf("%x", bs), nil
}
