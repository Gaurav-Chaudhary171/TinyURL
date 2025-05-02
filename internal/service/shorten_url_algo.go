package service

import (
	"crypto/sha256"
	"encoding/base64"
)

func ShortenURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8]
}
