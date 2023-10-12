package vaultclientparser

import (
	"errors"
	"strings"
)

var ErrUnknownVaultClient = errors.New("could not parse vault client")

// This exists to allow for unit testing of the parsing because creating a new
// client requires access to the respective vault
func ParseClient(clientstring string) (string, error) {
	clientstring = strings.ToLower(clientstring)

	if strings.Contains("hashicorp", clientstring) {
		return "Hashicorp", nil
	}

	return "", ErrUnknownVaultClient
}
