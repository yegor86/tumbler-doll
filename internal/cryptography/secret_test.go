package cryptography

import (
	"testing"
)

func Test_aes128gcm_encrypt_decrypt(t *testing.T) {

	secret := &Secret{}
	plainText := "SecretText"
	encryptionKey := []byte("fEfakgn@dsf#fgff")

	encrypted, err := secret.encryptAes128Gcm([]byte(plainText), encryptionKey)
	if err != nil {
		t.Fatalf("Failed to encrypt plain text: %v", err)
	}

	decrypted, err := secret.decryptAes128Gcm(encrypted, encryptionKey)
	if err != nil {
		t.Fatalf("Failed to decrypt cipher: %v", err)
	}
	if decrypted != plainText {
		t.Errorf("Expected '%s', got '%s'", plainText, decrypted)
	}
}
