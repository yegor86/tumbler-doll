/*
MIT License

Copyright (c) 2019 Andrzej Rehmann

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cryptography

import (
	"crypto/aes"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
)

const (
	magicChecksum = "::::MAGIC::::"
)

/*
  Decrypts hudson.util.Secret using the master.key
  1. master.key is hashed and trimmed to 16 bytes
  2. master key is used to decrypt hudson.util.Secret with AES-128-ECB
  3. decrypted secret is trimmed to 16 bytes
  4. secret is returned, later to be used for decrypting Jenkins credentials with AES-128-ECB
*/
func DecryptHudsonSecret(masterKey []byte, hudsonSecret []byte) ([]byte, error) {
	hashedMasterKey := hashMasterKey(masterKey)
	decryptedSecret := decryptAes128Ecb(hudsonSecret, hashedMasterKey)

	return decryptedSecret[:16], nil
	// if secretContainsChecksum(decryptedSecret) {
	// 	return decryptedSecret[:16], nil
	// } else {
	// 	return nil, createError(decryptedSecret)
	// }
}

func createError(decryptedSecret []byte) error {
	msg := fmt.Sprintf(
		"Error. Decrypted hudson secret does not contain expected checksum.\n"+
			"Expected checksum keyword:\n\t%s\n"+
			"Decrypted secret:\n\t%q",
		magicChecksum,
		decryptedSecret)
	return errors.New(msg)
}

func secretContainsChecksum(encryptedSecret []byte) bool {
	return strings.Contains(string(encryptedSecret), magicChecksum)
}

/*
   Hash needs to be 16 bytes as Jenkins uses AES-128 encryption.
*/
func hashMasterKey(masterKey []byte) []byte {
	hasher := sha256.New()
	hasher.Write(masterKey)
	return hasher.Sum(nil)[:16]
}

/*
   ECB mode is deprecated and not included in golang crypto library.
*/
func decryptAes128Ecb(encryptedData []byte, key []byte) []byte {
	cipher, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(encryptedData))
	// size := 16
	// for bs, be := 0, size; bs < len(encryptedData); bs, be = bs+size, be+size {
	// 	cipher.Decrypt(decrypted[bs:be], encryptedData[bs:be])
	// }
	for i := 0; i < len(encryptedData); i += aes.BlockSize {
		cipher.Decrypt(decrypted[i:i+aes.BlockSize], encryptedData[i:i+aes.BlockSize])
	}
	return decrypted
}