package handlers

import (
	"crypto/sha256"
	"encoding/base64"
)

// ShortenURL generates a shortened URL from the original URL
func ShortenURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8] // Return first 8 characters
}
