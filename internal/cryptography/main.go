package cryptography

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	customCryptoLib "github.com/yegor86/tumbler-doll/internal/jenkins/cryptography"
)

const (
	defaultFileMode = 0640
)

type Cryptography struct {
	masterKeyData []byte
	secretKeyData []byte
}

var (
	instance *Cryptography
	once     sync.Once
)

func GetInstance() *Cryptography {
	once.Do(func() {
		instance = &Cryptography{}
	})
	return instance
}

// LoadOrSeedCrypto load or seed encrypted files
func (crypto *Cryptography) LoadOrSeedCrypto() error {

	if os.Getenv("JENKINS_HOME") == "" {
		log.Fatalln("JENKINS_HOME environment variable must be initialized")
	}

	secretsPath := filepath.Join(os.Getenv("JENKINS_HOME"), "secrets")
	err := os.MkdirAll(secretsPath, 0740)
	if err != nil {
		log.Fatalf("failed to create '$JENKINS_HOME/secrets' directory: %v\n", err)
	}

	crypto.masterKeyData = loadOrSeed("master.key", func() []byte {
		return GenerateKey(256)
	})
	encryptedSecret := loadOrSeed("hudson.util.Secret", func() []byte {
		return customCryptoLib.EncryptHudsonSecret(crypto.masterKeyData, GenerateKey(256))
	})
	crypto.secretKeyData, err = customCryptoLib.DecryptHudsonSecret(crypto.masterKeyData, encryptedSecret)
	if err != nil {
		log.Fatalf("failed to decrypt hudson secret: %v\n", err)
	}

	return nil
}

func loadOrSeed(fileName string, keyGenerator func()[]byte) []byte {
	keyPath := filepath.Join(os.Getenv("JENKINS_HOME"), "secrets", fileName)
	
	_, err := os.Stat(keyPath)
	if err != nil && os.IsNotExist(err) {
		os.WriteFile(keyPath, keyGenerator(), defaultFileMode)
	} else if (err != nil) {
		log.Fatalf("error loading key: %v\n", err)
	}
	
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("error loading key: %v\n", err)
	}
	return keyData
}
