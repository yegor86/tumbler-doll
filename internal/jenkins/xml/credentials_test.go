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
package xml

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	gitlab = &Credential{
		Tags: map[string]string{
			"scope":       "GLOBAL",
			"id":          "cd8ad67a-30a3-4261-b56a-300b860f2764",
			"description": "Gitlab admin user",
			"username":    "gitlabadmin",
			"password":    "{AQAAABAAAAAgyJpDg7KJuiCXs6hdfhD8xmnkTQPmlwLqXioAHQYbgpwbLHgAr928te6rYAIEIJlO}",
		},
	}
	bastion = &Credential{
		Tags: map[string]string{
			"scope":            "GLOBAL",
			"id":               "d90e353f-0fcf-4ce5-aa53-c43235999251",
			"description":      "Production bastion ssh key",
			"username":         "root",
			"passphrase":       "{AQAAABAAAAAQYt6YziqpO92++ZbOA5Bbua9x5bqkM7qMxWCfeUMGMyc=}",
			"privateKeySource": "",
			"privateKey":       "{AQAAABAAAAOAU9kOrNdjrQdV8FG1pcYnT3uCmamX8qsPvcB3w9UWzEBjUtydMkiRFKka78Z08kVKOs/3jfKpWP1FdQlzPRhARA0+k7/0+jDkMcIUpcAYDpI0SJlTq6hxxLxd8+59QNeQGCzw5y6fezjpsvBx2/zXPDkgJbkmCHt+B2BlvD5oLI7O8mNxLccKHBvwXmLSevtnublgxtQ6B578eZ8L3GIAQAvAXTAtU1yg/akIvjPSRxfGa9md8SyedAMYzJidFQfGfJISmvSW+BS+78NEtpJNQdNZX9G8Cv2aQPfakRIOIwe/p17y6w/24reDNjjsfX58ODpJaezOTY10nMTuuYpaaOeEk2FXVfXCjiln22GgakHR3gY+Cvz+ZSlgN/QTGEhMHnBojOJYxUqWxD3ogo9gYwvez0mtsMo2H/6F+UM0BJg/U7qN3ru7pYAvjJW0FK7Mv5Jq6j59k7Riuui4AH6m8JCWwNbPZmxy87D3asrAH85Nn6egCl9bAD0Sf5BVNfeiPaB043LkRTlES33x9DopKG3dsSctkQO0XoH3cavgq66UNWikSao3XBLLmVQTGq38ZboEbPikep3mZHeVdafITcV258xhbi7Kw8+Qj3buJYfPReR9agKhyZGMjU1ZWuAShG43GinOjbRge3q9rA2isCk7icyh7pUeVWOS2R6RtMYH/dkNFcCJx0VgZ3FRxuEpLTEYmYfycTW/4lcWvFbeF9JqonDotaZiOPEDciG57fSRsIkp0uGD4iEYL6R8bJTM/tCfbmb3nbKkx5IcgjoZegbmEmK2RcQIOsV79DSo8HkDwLJplVPXCTKFD5soAp0EgszUK/xH8b2xCRAE0/5mPmzSwzOL+Dd2GNtJ9u/LhIneEtUZrPWPMRxH0Q4+4PVbiNv6sBMK+fLV7yJxl0TAl31qdfcwNF6bYm86ubywEVT/Og/SX2//Dn3ERI1ef3SCIrLMEWwNM/d3LhBRZbvBn9/Wm1qnfKBMVpvqyJC4RvkIsHNjpGEttNSYQ6RNQHIkZiJlZkLrKcdMrlz6i9SRzuWstRCg+nU22XSV4qYyIC3fOXuUuzyXiEvfwJlDAWT/y7HWKeXD8eQPACwtof2h9YGr/V+Kbl3L2lzxZwpu5n+DGhWmOThzZqF6scJVMQMdN840leC33lDDfFJHGk2KqzTLCr9Youk9N1lRWQpSJaYvMHcqZS0UsQ6XTjmFx+W+}",
		},
	}
)

func Test_reads_credentials_from_xml_file(t *testing.T) {
	expectedCredentials := []Credential{*gitlab, *bastion}
	credentialsXml, _ := os.ReadFile("../test/resources/credentials.xml")

	actualCredentials, _ := ParseCredentialsXml(credentialsXml)

	assert.Equal(t, expectedCredentials, actualCredentials)
}
