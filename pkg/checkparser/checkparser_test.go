package checkparser

import (
	"reflect"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/checks/datetime"
	"github.com/kudelskisecurity/youshallnotpass/pkg/checks/imagehash"
	"github.com/kudelskisecurity/youshallnotpass/pkg/checks/mfarequired"
	"github.com/kudelskisecurity/youshallnotpass/pkg/checks/scripthash"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
)

func TestParseChecks(t *testing.T) {
	parseChecksTests := []struct {
		name           string
		configs        []config.CheckConfig
		checkType      string
		ciPlatform     string
		expectedError  error
		expectedChecks []checks.Check
	}{
		{
			name: "test parse valid script hash check",
			configs: []config.CheckConfig{
				{
					Name: "scriptHash",
					Options: map[string]interface{}{
						"mfaOnFail": true,
					},
				},
			},
			checkType:     "script",
			ciPlatform:    "gitlab",
			expectedError: nil,
			expectedChecks: []checks.Check{
				&scripthash.ScriptHashCheck{},
			},
		},
		{
			name: "test parse invalid script hash check",
			configs: []config.CheckConfig{
				{
					Name: "scriptHash",
					Options: map[string]interface{}{
						"mfaOnFail": true,
					},
				},
			},
			checkType:      "image",
			ciPlatform:     "gitlab",
			expectedError:  nil,
			expectedChecks: []checks.Check{},
		},
		{
			name: "test parse valid image hash check",
			configs: []config.CheckConfig{
				{
					Name: "imageHash",
					Options: map[string]interface{}{
						"mfaOnFail": true,
					},
				},
			},
			checkType:     "image",
			ciPlatform:    "gitlab",
			expectedError: nil,
			expectedChecks: []checks.Check{
				&imagehash.ImageHashCheck{},
			},
		},
		{
			name: "test parse invalid image hash check",
			configs: []config.CheckConfig{
				{
					Name: "imageHash",
					Options: map[string]interface{}{
						"abortOnFail": true,
					},
				},
			},
			checkType:      "script",
			ciPlatform:     "gitlab",
			expectedError:  nil,
			expectedChecks: []checks.Check{},
		},
		{
			name: "test parse valid mfa required check",
			configs: []config.CheckConfig{
				{
					Name: "mfaRequired",
				},
			},
			checkType:     "all",
			ciPlatform:    "gitlab",
			expectedError: nil,
			expectedChecks: []checks.Check{
				&mfarequired.MfaRequiredCheck{},
			},
		},
		{
			name: "test parse valid date time check",
			configs: []config.CheckConfig{
				{
					Name: "dateTimeCheck",
				},
			},
			checkType:     "all",
			ciPlatform:    "gitlab",
			expectedError: nil,
			expectedChecks: []checks.Check{
				&datetime.DateTimeCheck{},
			},
		},
		{
			name: "test invalid test name",
			configs: []config.CheckConfig{
				{
					Name: "invalidCheck",
				},
			},
			checkType:      "all",
			ciPlatform:     "gitlab",
			expectedError:  ErrUnknownCheckNameError,
			expectedChecks: []checks.Check{},
		},
	}

	for _, test := range parseChecksTests {
		checks, err := ParseChecks(test.configs, "testJob", "", []string{""}, test.checkType, test.ciPlatform)

		if err != test.expectedError {
			t.Errorf("Unexpected error (%+v) in test: %s", err, test.name)
		}

		if len(checks) != len(test.expectedChecks) {
			t.Errorf("different checks length - expected: (%d) - got: (%d) - for: %s", len(test.expectedChecks), len(checks), test.name)
		}

		for i, check := range checks {
			if reflect.TypeOf(check) != reflect.TypeOf(test.expectedChecks[i]) {
				t.Errorf("different check types found - expected (%s) - got: (%s) - for %s", reflect.TypeOf(test.expectedChecks[i]), reflect.TypeOf(check), test.name)
			}
		}
	}
}
