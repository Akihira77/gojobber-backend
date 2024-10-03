package util

import (
	"math/rand"
	"mime/multipart"
	"time"
)

func RandomStr(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)

}

func ValidateImgExtension(file *multipart.FileHeader) bool {
	imgExtAllowed := []string{"image/webp", "image/png", "image/jpg", "image/jpeg"}

	for _, ext := range imgExtAllowed {
		if file.Header.Get("Content-Type") == ext {
			return true
		}
	}

	return false
}
