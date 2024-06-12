package mt

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

func CheckFileExist(filepath string) bool {
	_, error := os.Stat(filepath)
	return !os.IsNotExist(error)
}

func GetFileHashSHA256(fileName string) (string, error) {
	var returnSHA256String string

	file, err := os.Open(fileName)
	if err != nil {
		return returnSHA256String, err
	}
	defer file.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA256String, err
	}

	hashInBytes := hash.Sum(nil)
	returnSHA256String = hex.EncodeToString(hashInBytes)

	return strings.ToLower(returnSHA256String), nil
}
