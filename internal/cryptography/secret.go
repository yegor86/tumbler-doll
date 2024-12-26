package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type (
	Secret struct{
		Type string
	}
)

// GenerateKey generate a keyLength-character key using cryptographic random generator 
func GenerateKey(keyLength int) []byte {
	masterKey := make([]byte, keyLength)
	_, err := rand.Read(masterKey)
	if err != nil {
		panic(fmt.Errorf("error generating master key %w", err))
	}

	encodedKey := hex.EncodeToString(masterKey)
	return []byte(encodedKey)
}

// Encrypt encrypts plaintext using AES-GCM
func (secret *Secret) encryptAes128Gcm(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func (secret *Secret) decryptAes128Gcm(encryptedBase64 string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("invalid cipher text")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
