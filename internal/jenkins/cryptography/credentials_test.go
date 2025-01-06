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
	"testing"

	"github.com/yegor86/tumbler-doll/internal/jenkins/xml"
	"github.com/stretchr/testify/assert"
)

const (
	decryptedSecret = "V7sRJ]hBJE/7HWk4"
)

var (
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
				"password": "AQAAABAAAAAQTcjjf6+TSYQKfKW5GGFXONqrEhCYnQ1solOK6wQMfTc=",
			},
		},
		{
			Tags: map[string]string{
				"username": "gitlabadmin",
				"password": "{AQAAABAAAAAgyJpDg7KJuiCXs6hdfhD8xmnkTQPmlwLqXioAHQYbgpwbLHgAr928te6rYAIEIJlO}",
			},
		},
		{
			Tags: map[string]string{
				"username":   "root",
				"passphrase": "{AQAAABAAAAAQYt6YziqpO92++ZbOA5Bbua9x5bqkM7qMxWCfeUMGMyc=}",
				"privateKey": "{AQAAABAAAAOAU9kOrNdjrQdV8FG1pcYnT3uCmamX8qsPvcB3w9UWzEBjUtydMkiRFKka78Z08kVKOs/3jfKpWP1FdQlzPRhARA0+k7/0+jDkMcIUpcAYDpI0SJlTq6hxxLxd8+59QNeQGCzw5y6fezjpsvBx2/zXPDkgJbkmCHt+B2BlvD5oLI7O8mNxLccKHBvwXmLSevtnublgxtQ6B578eZ8L3GIAQAvAXTAtU1yg/akIvjPSRxfGa9md8SyedAMYzJidFQfGfJISmvSW+BS+78NEtpJNQdNZX9G8Cv2aQPfakRIOIwe/p17y6w/24reDNjjsfX58ODpJaezOTY10nMTuuYpaaOeEk2FXVfXCjiln22GgakHR3gY+Cvz+ZSlgN/QTGEhMHnBojOJYxUqWxD3ogo9gYwvez0mtsMo2H/6F+UM0BJg/U7qN3ru7pYAvjJW0FK7Mv5Jq6j59k7Riuui4AH6m8JCWwNbPZmxy87D3asrAH85Nn6egCl9bAD0Sf5BVNfeiPaB043LkRTlES33x9DopKG3dsSctkQO0XoH3cavgq66UNWikSao3XBLLmVQTGq38ZboEbPikep3mZHeVdafITcV258xhbi7Kw8+Qj3buJYfPReR9agKhyZGMjU1ZWuAShG43GinOjbRge3q9rA2isCk7icyh7pUeVWOS2R6RtMYH/dkNFcCJx0VgZ3FRxuEpLTEYmYfycTW/4lcWvFbeF9JqonDotaZiOPEDciG57fSRsIkp0uGD4iEYL6R8bJTM/tCfbmb3nbKkx5IcgjoZegbmEmK2RcQIOsV79DSo8HkDwLJplVPXCTKFD5soAp0EgszUK/xH8b2xCRAE0/5mPmzSwzOL+Dd2GNtJ9u/LhIneEtUZrPWPMRxH0Q4+4PVbiNv6sBMK+fLV7yJxl0TAl31qdfcwNF6bYm86ubywEVT/Og/SX2//Dn3ERI1ef3SCIrLMEWwNM/d3LhBRZbvBn9/Wm1qnfKBMVpvqyJC4RvkIsHNjpGEttNSYQ6RNQHIkZiJlZkLrKcdMrlz6i9SRzuWstRCg+nU22XSV4qYyIC3fOXuUuzyXiEvfwJlDAWT/y7HWKeXD8eQPACwtof2h9YGr/V+Kbl3L2lzxZwpu5n+DGhWmOThzZqF6scJVMQMdN840leC33lDDfFJHGk2KqzTLCr9Youk9N1lRWQpSJaYvMHcqZS0UsQ6XTjmFx+W+}",
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
