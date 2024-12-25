package cryptography

import (
	"testing"
)

func TestEncrypt(t *testing.T) {

	secret := &Secret{}
	plainText := "SecretText"
	encryptionKey := []byte("fEfakgn@dsf#fgff")

	encrypted, err := secret.Encrypt([]byte(plainText), encryptionKey)
	if err != nil {
		t.Fatalf("Failed to encrypt plain text: %v", err)
	}

	decrypted, err := secret.Decrypt(encrypted, encryptionKey)
	if err != nil {
		t.Fatalf("Failed to decrypt cipher: %v", err)
	}
	if decrypted != plainText {
		t.Errorf("Expected '%s', got '%s'", plainText, decrypted)
	}
}
