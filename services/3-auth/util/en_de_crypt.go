package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
)

func encrypt(plaintext string) ([]byte, error) {
	secretKey := os.Getenv("EN_DE_CRYPTION_KEY")
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		log.Printf("Error creating AES cipher block:\n%+v", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		log.Printf("Error setting up GCM mode:\n%+v", err)
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Reader.Read(nonce)
	if err != nil {
		log.Printf("Error generating nonce:\n%+v", err)
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return ciphertext, nil
}

func EncryptAndEncodeToHex(plaintext string) (string, error) {
	ciphertext, err := encrypt(plaintext)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), err
}

func decrypt(ciphertext string) (string, error) {
	secretKey := os.Getenv("EN_DE_CRYPTION_KEY")
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func DecodeToStringAndDecrypt(hextext string) (string, error) {
	hexBytes, err := hex.DecodeString(hextext)
	if err != nil {
		return "", err
	}

	plaintext, err := decrypt(string(hexBytes))
	if err != nil {
		return "", err
	}

	return plaintext, nil
}
