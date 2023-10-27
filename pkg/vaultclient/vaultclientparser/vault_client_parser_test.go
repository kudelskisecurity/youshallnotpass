package vaultclientparser

import (
	"testing"
)

func TestVaultClientParser(t *testing.T) {
	parseVaultClientTests := []struct {
		name           string
		clientName     string
		expectedClient string
		error          error
	}{
		{
			name:           "test parse hashicorp from lowercase",
			clientName:     "hashicorp",
			expectedClient: "Hashicorp",
			error:          nil,
		},
		{
			name:           "test parse hashicorp from uppercase",
			clientName:     "HASHICORP",
			expectedClient: "Hashicorp",
			error:          nil,
		},
		{
			name:           "test parse hashicorp from odd capitalization",
			clientName:     "HaShIcOrP",
			expectedClient: "Hashicorp",
			error:          nil,
		},
		{
			name:           "test parse unknown client",
			clientName:     "ashjfaoiuehpao",
			expectedClient: "",
			error:          ErrUnknownVaultClient,
		},
	}

	for testNum, test := range parseVaultClientTests {
		client, err := ParseClient(test.clientName)

		if err != test.error {
			t.Errorf("\n%d) unexpected error occurred - \nexpected: (%+v)\ngot: (%+v)", testNum, test.error, err)
		}

		if client != test.expectedClient {
			t.Errorf("\n%d) unexpected client parsed - \nexpected: (%s)\ngot: (%s)", testNum, test.expectedClient, client)
		}
	}
}
