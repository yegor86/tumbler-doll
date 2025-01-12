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
	"github.com/yegor86/tumbler-doll/internal/jenkins/xml"
)

const (
	decryptedSecret = "V7sRJ]hBJE/7HWk4"
)

var (
	// Initialization vector for encryption
	iv = []byte{77, 200, 227, 127, 175, 147, 73, 132, 10, 124, 165, 185, 24, 97, 87, 56}

	oldFormatEncryptedCredentials = []xml.Credential{
		{
			Tags: map[string]string{
				"username": "user_1",
				"password": "jteGsiX7VHD7320kqlnhY3doAdfZx5VoEmG+VrkkMco=",
			},
		},
		{
			Tags: map[string]string{
				"username":   "user_1",
				"privateKey": "LINE_1\nLINE_2\nLINE_2",
				"passphrase": "Z7N3IHxxJIUpFQWPbDdZz32cuRUtc2goojPSub0AMvw=",
			},
		},
	}

	oldFormatDecryptedCredentials = []xml.Credential{
		{
			Tags: map[string]string{
				"username": "user_1",
				"password": ")gYoLsi{D0A[#!shE}iw4it}c"},
		},
		{
			Tags: map[string]string{
				"passphrase": "CZ)h'-z*+&cTUL)5n9P:A&7t'bs'q",
				"username":   "user_1",
				"privateKey": "LINE_1\nLINE_2\nLINE_2"},
		},
	}

	newFormatEncryptedCredentials = []xml.Credential{
		{
			Tags: map[string]string{
				"username": "xfireadmin",
				"password": "{AQAAABAAAAAQTcjjf6+TSYQKfKW5GGFXONqrEhCYnQ1solOK6wQMfTc=}",
			},
		},
		{
			Tags: map[string]string{
				"username": "gitlabadmin",
				"password": "{AQAAABAAAAAgTcjjf6+TSYQKfKW5GGFXOF6x26+HXozYUgVbhYE7Xob3Z52sEAJRbydS3dn2yoL3}",
			},
		},
		{
			Tags: map[string]string{
				"username":   "root",
				"passphrase": "{AQAAABAAAAAQTcjjf6+TSYQKfKW5GGFXODUpsDKtsf3FhdpoNK1PC2I=}",
				"privateKey": "{AQAAABAAAAOATcjjf6+TSYQKfKW5GGFXOD2w2E2B6Nn3nThLPLC2L8Fvg4MgY4Rkx9gn17nTelvSMJ3/YKozZrrmXQ9C0IJ67aA0F4ptCosn9UYxg6WfHZnQcqud+JtYBAs/sH3/eGPHhKuZmVqDkhy9U2S8b7NimLSHvRGaAl0NcIIZGoeXm9b8OrokF2YS00KWbmfW0wARsQDSUOsFaxCdEe1pUBTTyPf8MBr6WJZ0DDnpW5jVF6DUoJCwJPX1PyGntisu5y8DsbDguhjsYGGl0ojh1Shkfwvrphjz4n+lOZFiFl81eOV6wT2cyoChPCPGsjV+eHz9DOh5Br1noQ1HhQryD5rWT0EahYJFTOjQNbPPMwqFO9QIjB9ME/UvZuVzGb+cx/46eWASmatTcxpbsZKGMxgWodSvby0SzI0T/nlAlCKWzALtGF52JHK1KGLjrbrtX8yJltZqqqridWJfF+kJbmcAA4mR31IpMtt+pWd+8usWWx+AZrfHP/UBWb6pbndUhR6ew6645PiUe+6gwc9AoWvvkG92naUr7JQTtRBAuheVdpgSq8Haco6UPnf/+RyVe37qLJYmBtT33WDsRVy969B6uhARO468mzeM2t9mLiRBq7/Tvg9IIUh2GhevV+pYAVcnIPVKf4ElW1ZOrTpUs2xMqTmB3YuttwHp2gnifPnubbO1qvbbCBW+0KT+fVZSRgNBZ1+kPGz7MBJAK82CmDP1bxe4ulX7nG4ZlVtCQykfhZeYckyBNVtQU5j7ER2TY833Zi1jjIsdgoULlSOat8VDYZLqoTnVmrIzDoMbK48SL1n0YvjochXNjZoP1cDW57kcg0KlgvPTwYHpz3MDaHzQ8l0BblPn3Px6sCybZFEflBj5yjFAvl7twKkNRbHJvdf+aIxwUOVSzO1vABpqc/Zd8ASHRtHyUTlXOYy49RjfHq2iI/aFScYQpTpvNj6A3lm6GSl5ARMhBV/Jmnp2X5uwwyLGUY7/Ptd8tnw4ruJZCLtXaHvBC3IdHrmX0W0iIPcIA33fVx0NFyUoNH1HGRja8dmwMbwslc6o3nsbEMBEAD0FjA5tihJ6EBzRLxiH4nqcAJGQZoRllhi76I1aSavkgl2u6jhVPfIJSwgvjEKueHyoxB3yrXKQ1qth0bzi6boIry73wdYx5zKlFdvLhPOTeeEm+9fkz4Yr0TnNu/uf295EJLpyfak5f3Z97Z+wFovi}",
			},
		},
	}
	newFormatDecryptedCredentials = []xml.Credential{
		{
			Tags: map[string]string{
				"username": "xfireadmin",
				"password": "ilovexfire",
			},
		},
		{
			Tags: map[string]string{
				"username": "gitlabadmin",
				"password": "Drmhze6EPcv0fN_81Bj",
			},
		},
		{
			Tags: map[string]string{
				"username":   "root",
				"passphrase": "IEPkO5nYhEG",
				"privateKey": `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`,
			},
		},
	}
)

func Test_decrypts_old_format_credentials(t *testing.T) {
	secret := []byte(decryptedSecret)

	credentials, _ := DecryptCredentials(oldFormatEncryptedCredentials, secret)

	assert.Equal(t, credentials, oldFormatDecryptedCredentials)
}

func Test_decrypts_new_format_credentials(t *testing.T) {
	secret := []byte(decryptedSecret)

	credentials, _ := DecryptCredentials(newFormatEncryptedCredentials, secret)

	assert.Equal(t, credentials, newFormatDecryptedCredentials)
}

func Test_encrypts_new_format_credentials(t *testing.T) {
	
	secret := []byte(decryptedSecret)

	encrypt := func(plaintext, key []byte) ([]byte, error) {
		return _encryptAes128Cbc(plaintext, key, iv)
	}
	actualCredentials, err := EncryptCredentials(newFormatDecryptedCredentials, secret, encrypt)
	if err != nil {
		t.Fatalf("Failed to encrypt credentials xml: %v", err)
	}

	assert.Equal(t, newFormatEncryptedCredentials, actualCredentials)
}

func Test_decrypt_credentials_from_xml_file(t *testing.T) {
	t.Skip("skipping testing")

	expectedCredentials := []xml.Credential{}
	credentialsXml, _ := os.ReadFile("../test/resources/credentials.xml")

	credentials, _ := xml.ParseCredentialsXml(credentialsXml)

	encrypt := func(plaintext, key []byte) ([]byte, error) {
		return _encryptAes128Cbc(plaintext, key, iv)
	}
	encryptedCredentials, _ := EncryptCredentials(credentials, []byte(decryptedSecret), encrypt)

	os.WriteFile("../test/resources/encrypted_credentials.xml", []byte(encryptedCredentials[2].Tags["privateKey"]), 0770)

	assert.Equal(t, expectedCredentials, encryptedCredentials)
}
