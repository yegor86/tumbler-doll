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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// hudson.util.Secret contains '::::MAGIC::::' checksum when it's created by Jenkins
	encryptedHudsonSecretWithNoMagic = "UKFS\xe6\xedv\xf8\x987xo^\x83\xf2*\xb4\x03\x97.\xe7\xd2\xde\x14\xa3\xb6\xcfF\x9c\xa3^)q\xa5\xa0\x85h\xb1'\xaa\xb4\xad\xdb\xf0\x8dKe\x06\xa0k˪ˣX\x8c6\xe6\x14V\x86:\xeb\x1dq/\xfa\xaf?\xf5®>\xec\x83HC\x83\xf9\xc2\xf1qo\x87\x9f\xef(\xed\x06\xb7\x93`\xf2fC\xccy\xe6\xe0Bۙ\xcc\x1e5[\x9c\x9b\xa0K\x9e\xab\xb09\xecA\x1d19HS\xb3<\xa7\xa4\xec\xce\xf3\xb7\xe5\xde\x10J\x06\xdeK9zj\x85\x95\xcee\x19=}W-h\xfb\xb9\x121V\xb1F\xeeK\xf1\t\xe8\x87\xf6d\xe1\xb8\xfd%\xca:a\xdcnH\xdf\xfc\xd2\xc9[\xf8e-΅\xab\xbc\x04\xdfK'1j%\xbe\x93\x12\xfb\x00\x8a\x89\x84\xc1\x1f`\x9bڏy\xedMc\xfcGrh\xcf\x1e\xef!~\xec\xbd\xf5\xba\x97]u\xffr\t\xf7\x19X9\xcfo\xce\x15}l\xbaM\x89~\xe5s\xed\xd8:\xb6ᓋRX\x84#\xabu[\x07\xf8\xde\x1awH\xc2;b \x04\xc3"

	hudsonSecret = "V7sRJ]hBJE/7HWk4vr=TkdJ38T+siDd"
)

func Test_aes128cbc_encryption(t *testing.T) {
	secret := hudsonSecret[:16]
	plaintext := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCwgxo7cl2RajAWFseL0JAIBJbZ6dFWBGcq7+TMkP8viDwfLj4u
iYqERw+Y/lW0VZxuQuVMBfcCCINTG0S3W+MYPKiHKSaQWV53oOyPUCWaU1WjMHG4
Y3DFeE8NomOqLEOjCHDAIkZDzeEO14S8OW0fEycvR7Opo8lI6TJ4xAi9gwIDAQAB
AoGAJDu1TdCrLmd62X3xllzIxCyU/sSFiT+8Ic8+y1NUXuB7XvcyIoFvYrnnlMNY
unz8cJHg2ds7mjo/IvctAuqk0gQ5cMUgxf5QzP0ZlIcHyq8lB0YsRki0aZHSuv7r
N2KKUrayPTAJrA08GPYA+koLnY1R/yNiauo3D0ERe7KdCUECQQDolVFRXQ3KAlsu
ZGpMXbjm3H+V610D1/F4xg8qae5tKZTpbPwSnhHvCYLfNE6CSGf4JE/nhhWQk5th
j8oB4A0hAkEAwkiWXgckr+w2q+nApKDu5co27vdOPmC0q+8gk1+hCLs6Tn5V0dVx
z3aqUdQpiIYOrd5FSfXa7YNnO8VgmoayIwJAK/XFB+7ho1PsrgkWulZgk2oLx2dU
DlzrbBtrVGXvRby9Q51wy4gK9bZDgTKewCs1U4Zxf94tB0WO8dK+qLoTYQJBALcQ
4KcfAgHGiUl6C+zUO+dIoHSRkSeTxgpQW5iiPkHU8b7uqfz7q676OMi8Kpqa/w/z
5cQoJq8w50BZ3oocq5MCQBew/PwOfusahnBiUoFY0CfWTR4HZ86Uo1zgtPKoLCUG
hDA6SHkmIEPkO5nYhEGMryddRI7rsB4EKJaQ8AnJ7r4=
-----END RSA PRIVATE KEY-----`)

	encryptedText, _ := encryptAes128Cbc(plaintext, []byte(secret))

	decryptedText, _ := decryptAes128Cbc(encryptedText, []byte(secret))

	assert.Equal(t, plaintext, decryptedText)
	assert.True(t, len(decryptedText) > 1)
}

func Test_encrypts_secret(t *testing.T) {
	masterKey, _ := os.ReadFile("../test/resources/master.key")

	encryptedHudsonSecret := EncryptHudsonSecret(masterKey, []byte(hudsonSecret))

	actualDecryptedHudsonSecret, _ := DecryptHudsonSecret(masterKey, encryptedHudsonSecret)

	assert.Equal(t, []byte(hudsonSecret), actualDecryptedHudsonSecret)
	assert.True(t, len(actualDecryptedHudsonSecret) > 1)
}

func Test_decrypts_secret(t *testing.T) {
	masterKey, _ := os.ReadFile("../test/resources/master.key")
 	encryptedSecret, _ := os.ReadFile("../test/resources/hudson.util.Secret")
	// encryptedSecret := []byte(encryptedHudsonSecretWithNoMagic)
	expectedDecryptedHudsonSecret := []byte(hudsonSecret)

	actualDecryptedHudsonSecret, _ := DecryptHudsonSecret(masterKey, encryptedSecret)

	assert.Equal(t, expectedDecryptedHudsonSecret, actualDecryptedHudsonSecret)
	assert.True(t, len(actualDecryptedHudsonSecret) > 1)
}
