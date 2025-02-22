package util

import (
	"crypto/sha256"
	"fmt"
)

func SanitizeString(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))[:8]
}

func Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:])
}
