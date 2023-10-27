package vaultclient

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateRandomBytes(t *testing.T) {
	generateRandomBytesTests := []struct {
		name  string
		bytes int
		error error
	}{
		{
			name:  "generate valid random bytes",
			bytes: 13,
			error: nil,
		},
		{
			name:  "generate random bytes negative",
			bytes: -20,
			error: ErrUnsupportedNumBytes,
		},
	}

	for testNum, test := range generateRandomBytesTests {
		bytes, err := generateRandomBytes(test.bytes)
		if err != nil {
			if err != test.error {
				t.Errorf("\n%d) unexpected error -\nexpected: (%+v)\ngot: (%+v)", testNum, test.error, err)
			}
		} else if len(bytes) != test.bytes {
			t.Errorf("\n%d) unexpected number of bytes generated -\nexpected: (%d)\ngot: (%d)", testNum, test.bytes, len(bytes))
		}
	}
}

func TestGenerateRandomStringURLSafe(t *testing.T) {
	generateRandomStringTests := []struct {
		name      string
		bytes     int
		urlLength int
		error     error
	}{
		{
			name:      "generate valid random string safe",
			bytes:     13,
			urlLength: 18,
			error:     nil,
		},
		{
			name:      "generate invalid random string negative bytes",
			bytes:     -20,
			urlLength: 0,
			error:     ErrUnsupportedNumBytes,
		},
	}

	for testNum, test := range generateRandomStringTests {
		str, err := GenerateRandomStringURLSafe(test.bytes)
		if test.error != err {
			t.Errorf("\n%d) unexpected error -\nexpected: (%+v)\ngot: (%+v)", testNum, test.error, err)
		}

		if len(str) != test.urlLength {
			t.Errorf("\n%d) unexpected url length -\nexpected: (%d)\ngot: (%d)", testNum, test.urlLength, len(str))
		}
	}
}

func TestCreateValidationToken(t *testing.T) {
	generateRandomStringIgnoreError := func() string {
		str, _ := GenerateRandomStringURLSafe(12)
		return str
	}

	createValidationTokenTests := []struct {
		name      string
		tokenName string
	}{
		{
			name:      "create validation token named validation_token",
			tokenName: "validation_token",
		},
		{
			name:      "create random validation token",
			tokenName: generateRandomStringIgnoreError(),
		},
	}

	for testNum, test := range createValidationTokenTests {
		err := CreateValidationToken(test.tokenName)
		if err != nil {
			t.Errorf("\n%d) Expected err to be nil, but got %s", testNum, err.Error())
			return
		}

		dir, err := os.Getwd()
		if err != nil {
			t.Errorf("\n%d) Expected err to be nil, but got %s", testNum, err.Error())
			return
		}

		_, err = os.Stat(dir + "/" + test.tokenName)
		if err != nil {
			t.Errorf("\n%d) File was not created successfully", testNum)
			return
		}

		err = os.Remove(dir + "/" + test.tokenName)
		if err != nil {
			t.Errorf("\n%d) Could not remove the validation token file", testNum)
			return
		}
	}
}

func TestValidationTokenCreated(t *testing.T) {
	err := CreateValidationToken("validation_token")
	if err != nil {
		t.Errorf("Expected err to be nil, but got %s", err.Error())
		return
	}

	created := TokenCreated("validation_token")
	if !created {
		t.Error("Expected the token to exist")
		return
	}

	dir, _ := os.Getwd()
	err = os.Remove(dir + "/" + "validation_token")
	if err != nil {
		t.Error("Could not remove the validation token file")
		return
	}
}

func TestCreateValidationTokenTwice(t *testing.T) {
	err := CreateValidationToken("validation_token")
	if err != nil {
		t.Errorf("Expected err to be nil, but got %s", err.Error())
		return
	}

	err = CreateValidationToken("validation_token")
	print(err.Error())
	if !strings.Contains(err.Error(), "pre-validation file already exists") {
		t.Error("Expected pre-validation file to already exist")
		return
	}

	dir, _ := os.Getwd()
	err = os.Remove(dir + "/" + "validation_token")
	if err != nil {
		t.Error("Could not remove the validation token file")
		return
	}
}
