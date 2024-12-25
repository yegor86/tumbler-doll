/*
MIT License

# Copyright (c) 2019 Andrzej Rehmann

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
	"bytes"
	"crypto/aes"
	cipherLib "crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
	decryptedSecret, err := decryptAes128Ecb(hudsonSecret, hashedMasterKey)
	if err != nil {
		return nil, err
	}
	
	return decryptedSecret[:16], nil
	// if secretContainsChecksum(decryptedSecret) {
	// 	return decryptedSecret[:16], nil
	// } else {
	// 	return nil, createError(decryptedSecret)
	// }
}

func EncryptHudsonSecret(masterKey []byte, hudsonSecret []byte) ([]byte, error) {
	hashedMasterKey := hashMasterKey(masterKey)
	hudsonSecret = append(hudsonSecret, []byte(magicChecksum)...)
	decryptedSecret := encryptAes128Ecb(hudsonSecret, hashedMasterKey)
	
	return decryptedSecret, nil
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


// Hash needs to be 16 bytes as Jenkins uses AES-128 encryption.
func hashMasterKey(masterKey []byte) []byte {
	hasher := sha256.New()
	hasher.Write(masterKey)
	return hasher.Sum(nil)[:16]
}

// AES128 ECB mode
func decryptAes128Ecb(encryptedData []byte, key []byte) ([]byte, error) {
	cipher, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(encryptedData))
	size := 16
	for bs, be := 0, size; bs < len(encryptedData); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], encryptedData[bs:be])
	}
	
	return removePadding(decrypted)
}

// AES128 CBC mode
func decryptAes128Cbc(cipher []byte, secret []byte) ([]byte, error) {
	ivLength := binary.BigEndian.Uint32(cipher[1:5])
	dataLength := int(binary.BigEndian.Uint32(cipher[5:9]))

	cipher = cipher[9:] // strip version, iv and data length

	iv := cipher[:ivLength]
	cipher = cipher[ivLength:] //strip iv
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	mode := cipherLib.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipher, cipher)

	
	// Remove PKCS padding
	n := int(cipher[dataLength-1])
	return cipher[:dataLength-n], nil
}

// encryptAes128Ecb encrypts plaintext using AES-128 in ECB mode.
func encryptAes128Ecb(plaintext []byte, key []byte) []byte {
	// Create AES cipher block
	cipher, _ := aes.NewCipher(key)

	// Pad plaintext to match block size
	plaintext = padToBlockSize(plaintext, aes.BlockSize)

	// Encrypt each block
	encrypted := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += aes.BlockSize {
		cipher.Encrypt(encrypted[i:i+aes.BlockSize], plaintext[i:i+aes.BlockSize])
	}

	return encrypted
}

// encryptAes128Cbc encrypts the plaintext using AES-128 in CBC mode.
func encryptAes128Cbc(plaintext, key []byte) ([]byte, error) {
    if len(key) != 16 {
        return nil, errors.New("key length must be 16 bytes for AES-128")
    }

    cipher, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    plaintext = padToBlockSize(plaintext, aes.BlockSize)

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Encrypt the plaintext
	ciphertext := make([]byte, len(plaintext))
	mode := cipherLib.NewCBCEncrypter(cipher, iv)
	mode.CryptBlocks(ciphertext, plaintext)

    // Create a buffer to hold the output
	var output bytes.Buffer

	output.WriteByte(0x1)

	// Encode the IV length (4 bytes)
	ivLength := uint32(len(iv))
	if err := binary.Write(&output, binary.BigEndian, ivLength); err != nil {
		return nil, err
	}

	// Encode the data length (4 bytes)
	dataLength := uint32(len(ciphertext))
	if err := binary.Write(&output, binary.BigEndian, dataLength); err != nil {
		return nil, err
	}

	// Write the IV
	output.Write(iv)

	// Write the ciphertext
	output.Write(ciphertext)

    return output.Bytes(), nil
}

// padToBlockSize applies PKCS#7 padding to make the data a multiple of the block size.
func padToBlockSize(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// removePadding removes PKCS padding from decrypted data.
func removePadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	// The value of the last byte is the padding length
	padding := int(data[len(data)-1])

	// Ensure the padding value is valid
	if padding <= 0 || padding > len(data) {
		return nil, errors.New("invalid padding")
	}

	// Check if all the padded bytes have the correct value
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding bytes")
		}
	}

	// Return the data without padding
	return data[:len(data)-padding], nil
}