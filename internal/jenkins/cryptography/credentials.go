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
	"encoding/base64"
	"regexp"
	"strings"

	"github.com/yegor86/tumbler-doll/internal/jenkins/xml"
)

/*
  This is some next level reverse engineering.
  Kudos to http://xn--thibaud-dya.fr/jenkins_credentials.html
*/
func DecryptCredentials(credentials []xml.Credential, secret []byte) ([]xml.Credential, error) {
	decryptedCredentials := make([]xml.Credential, len(credentials))
	copy(decryptedCredentials, credentials)

	for i, credential := range credentials {
		for key, value := range credential.Tags {
			if isBase64EncodedSecret(value) {
				encodedCipher := stripBrackets(value)
				cipher, err := base64Decode(encodedCipher)
				if err != nil {
					return nil, err
				}

				decrypted, err := decrypt(cipher, secret)
				if err != nil {
					return nil, err
				}
				decryptedCredentials[i].Tags[key] = string(decrypted)
			}
		}

	}
	return decryptedCredentials, nil
}

/*
  New format of declaring a field to be a "base64 decoded secret" is by using {} brackets.
  Example:

    <password>{AQAAABAAAAAgPT7JbBVgyWiivobt0CJEduLyP0lB3uyTj+D5WBvVk6jyG6BQFPYGN4Z3VJN2JLDm}</password>

  Old format does not use the {} brackets.
  Instead jenkins seems to be usually suffixing the encoding with '=' sign.
  Example:

     <password>B+4pJjkJXD+pzyT9lcq8M8vF+p5YU4HmWy+MWldEdG4=</password>

  I'm not sure how to distinguish other encoded secrets from the "old days of jenkins".
  I don't want to comprehend Jenkins code from 4 years ago just to handle some edge cases.
  I can't try to decode all values because there are some phrases which
  would be false positive e.g. "root" (which is a valid base64 encoding)
*/
func isBase64EncodedSecret(text string) bool {
	if isBracketed(text) {
		encoded := textBetweenBrackets(text)
		return isBase64Encoded(encoded)
	}
	if strings.HasSuffix(text, "=") {
		return isBase64Encoded(text)
	}
	return false
}

func isBase64Encoded(text string) bool {
	_, err := base64.StdEncoding.DecodeString(text)
	if err == nil {
		return true
	}
	return false
}

func stripBrackets(text string) string {
	if isBracketed(text) {
		return textBetweenBrackets(text)
	}
	return text
}

func isBracketed(text string) bool {
	return strings.HasPrefix(text, "{") && strings.HasSuffix(text, "}")
}

func base64Decode(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func textBetweenBrackets(text string) string {
	return regexp.MustCompile("{(.*?)}").FindStringSubmatch(text)[1]
}

func decrypt(cipher []byte, secret []byte) ([]byte, error) {
	if cipher[0] == 1 {
		return decryptAes128Cbc(cipher, secret)
	}
	return decryptAes128Ecb(cipher, secret)
}
