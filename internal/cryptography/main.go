package cryptography

import (
	"fmt"
	"os"
	"path/filepath"
)

// https://github.com/hoto/jenkins-credentials-decryptor
func InitCrypto() error {

	secretsPath := filepath.Join(os.Getenv("JENKINS_HOME"), "secrets")
	err := os.MkdirAll(secretsPath, 0700)
	if err != nil {
		return fmt.Errorf("failed to create '$JENKINS_HOME/secrets' directory: %w", err)
	}

	masterKeyPath := filepath.Join(os.Getenv("JENKINS_HOME"), "secrets", "master.key")
	_, err = os.Stat(masterKeyPath)
	if err != nil {
		return err
	}



	return nil
}
