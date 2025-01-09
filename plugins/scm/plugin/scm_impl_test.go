package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/yegor86/tumbler-doll/internal/cryptography"
	"github.com/yegor86/tumbler-doll/internal/jenkins/xml"
)

func Test_checkout_ssh_auth(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to init test: %v", err)
	}
	os.Setenv("WORKSPACE", filepath.Join(homeDir, "workspace"))

	scm := &ScmPluginImpl{
		logger: hclog.Default(),
	}

	crypto := cryptography.GetInstance()
	crypto.Credentials = []xml.Credential{
		{
			Tags: map[string]string{
				"id":       "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
				"username": "xfireadmin",
				"password": "ilovexfire",
			},
		},
		{
			Tags: map[string]string{
				"id":       "d90e353f-0fcf-4ce5-aa53-c43235999251",
				"username": "gitlabadmin",
				"password": "Drmhze6EPcv0fN_81Bj",
			},
		},
		{
			Tags: map[string]string{
				"id":         "12345-1234-4696-af25-123455",
				"username":   "git",
				"passphrase": "",
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

	args := map[string]interface{}{
		"branch":        "main",
		"credentialsId": "12345-1234-4696-af25-123455",
		"url":           "git@github.com:yegor86/tumbler-doll.git",
	}

	result, err := scm.Checkout(args)
	if err != nil {
		t.Fatalf("Failed to checkout repo: %v", err)
	}
	fmt.Println(result)
}
