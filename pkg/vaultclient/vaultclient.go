package vaultclient

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/loggerclient"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

var (
	ErrUnsupportedNumBytes = errors.New("unsupported number of bytes")
	ErrCannotCreateToken   = errors.New("unable to create validation token")
	ErrPrevalidationExists = errors.New("pre-validation file already exists")
)

type VaultClient interface {
	GetNamespaceConfig(string) (config.NamespaceConfig, error)
	GetProjectConfig(string) (config.ProjectConfig, error)
	ReadWhitelists(string, string) (whitelist.Whitelist, error)
	WriteScratch(int, string) (string, error)
	LogMFAInstructions(string, loggerclient.LoggerClient)
	WaitForMFA(int, string) bool
	Cleanup(bool, string, string) error
}

func generateRandomBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrUnsupportedNumBytes
	}

	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomStringURLSafe(n int) (string, error) {
	if n < 0 {
		return "", ErrUnsupportedNumBytes
	}

	b, err := generateRandomBytes(n)
	b64String := base64.URLEncoding.EncodeToString(b)
	string := strings.Trim(b64String, "=")

	return string, err
}

func TokenCreated(validationToken string) bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}

	_, err = os.Stat(dir + "/" + validationToken)

	return !os.IsNotExist(err)
}

func CreateValidationToken(validationToken string) error {
	dir, _ := os.Getwd()

	if !TokenCreated(validationToken) {
		file, err := os.Create(dir + "/" + validationToken)
		if err != nil {
			return ErrCannotCreateToken
		}

		defer func() {
			if err := file.Close(); err != nil {
				fmt.Printf("Error closing file: %s\n", err)
			}
		}()

		return nil
	}

	return ErrPrevalidationExists
}
