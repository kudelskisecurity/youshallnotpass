package scripthash

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

func TestNewScriptHashCheck(t *testing.T) {
	newScriptHashCheckTests := []struct {
		name          string
		config        config.CheckConfig
		expectedCheck ScriptHashCheck
	}{
		{
			name: "test abortOnFail true and mfaOnFail true",
			config: config.CheckConfig{
				Name: "scriptHash",
				Options: map[string]interface{}{
					"abortOnFail": true,
					"mfaOnFail":   true,
				},
			},
			expectedCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   true,
				scriptLines: []string{""},
			},
		},
		{
			name: "test abortOnFail false and mfaOnFail false",
			config: config.CheckConfig{
				Name: "scriptHash",
				Options: map[string]interface{}{
					"abortOnFail": false,
					"mfaOnFail":   false,
				},
			},
			expectedCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: false,
				mfaOnFail:   false,
				scriptLines: []string{""},
			},
		},
		{
			name: "test abortOnFail missing and mfaOnFail missing",
			config: config.CheckConfig{
				Name:    "scriptHash",
				Options: map[string]interface{}{},
			},
			expectedCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: false,
				mfaOnFail:   false,
				scriptLines: []string{""},
			},
		},
	}

	jobName := "testJob"
	scriptLines := []string{""}
	for testNum, test := range newScriptHashCheckTests {
		scriptHashCheck := NewScriptHashCheck(test.config, jobName, scriptLines)
		if !reflect.DeepEqual(scriptHashCheck, test.expectedCheck) {
			t.Errorf("%d) Checks not equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedCheck, scriptHashCheck)
		}
	}
}
func TestScriptHashCheck(t *testing.T) {
	scriptHashCheckTests := []struct {
		name            string
		scriptHashCheck ScriptHashCheck
		whitelist       whitelist.Whitelist
		expectedResult  checks.CheckResult
	}{
		{
			name: "test script in whitelist - success",
			scriptHashCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   true,
				scriptLines: []string{
					`$ echo "this script is for testing"`,
					`$ curl -X POST -H yourdomain.com"`,
				},
			},
			whitelist: whitelist.Whitelist{
				AllowedScripts: []string{
					"testJob@sha256:v7gXKQ__R6c76dShfVJoe8gvPOc-TJxHhv3QDdRp5kw=",
				},
			},
			expectedResult: checks.CheckResult{
				Name:    ScriptHashCheckName,
				Version: ScriptHashCheckVersion,
				Error:   nil,
				Abort:   false,
				Mfa:     false,
				Details: ScriptHashCheckSuccessDetails,
			},
		},
		{
			name: "test script not in (empty) whitelist - abort",
			scriptHashCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   false,
				scriptLines: []string{
					`$ echo "this script is for testing"`,
					`$ curl -X POST -H yourdomain.com"`,
				},
			},
			whitelist: whitelist.Whitelist{},
			expectedResult: checks.CheckResult{
				Name:    ScriptHashCheckName,
				Version: ScriptHashCheckVersion,
				Error:   nil,
				Abort:   true,
				Mfa:     false,
				Details: fmt.Sprintf(ScriptHashCheckAbortScriptDetails, "testJob", "sha256:v7gXKQ__R6c76dShfVJoe8gvPOc-TJxHhv3QDdRp5kw="),
			},
		},
		{
			name: "test script not in whitelist (updated) - abort",
			scriptHashCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   false,
				scriptLines: []string{
					`$ echo "this script is for testing"`,
					`$ curl -X POST -H yourdomain.com"`,
				},
			},
			whitelist: whitelist.Whitelist{
				AllowedScripts: []string{
					"testJob@sha256:v7gXKQ__R6c76dShfVJoe8gvPOc-TJxHhv3QDdRpxyz=",
				},
			},
			expectedResult: checks.CheckResult{
				Name:    ScriptHashCheckName,
				Version: ScriptHashCheckVersion,
				Error:   nil,
				Abort:   true,
				Mfa:     false,
				Details: fmt.Sprintf(ScriptHashCheckAbortScriptDetails+ScriptHashCheckUpdatedScriptDetails,
					"testJob", "sha256:v7gXKQ__R6c76dShfVJoe8gvPOc-TJxHhv3QDdRp5kw=", "testJob"),
			},
		},
		{
			name: "test script not in (empty) whitelist - mfa",
			scriptHashCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: false,
				mfaOnFail:   true,
				scriptLines: []string{
					`$ echo "this script is for testing"`,
					`$ curl -X POST -H yourdomain.com"`,
				},
			},
			whitelist: whitelist.Whitelist{},
			expectedResult: checks.CheckResult{
				Name:    ScriptHashCheckName,
				Version: ScriptHashCheckVersion,
				Error:   nil,
				Abort:   false,
				Mfa:     true,
				Details: fmt.Sprintf(ScriptHashCheckMfaScriptDetails, "testJob", "sha256:v7gXKQ__R6c76dShfVJoe8gvPOc-TJxHhv3QDdRp5kw="),
			},
		},
		{
			name: "test script not in whitelist (updated) - mfa",
			scriptHashCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: false,
				mfaOnFail:   true,
				scriptLines: []string{
					`$ echo "this script is for testing"`,
					`$ curl -X POST -H yourdomain.com"`,
				},
			},
			whitelist: whitelist.Whitelist{
				AllowedScripts: []string{
					"testJob@sha256:v7gXKQ__R6c76dShfVJoe8gvPOc-TJxHhv3QDdRpxyz=",
				},
			},
			expectedResult: checks.CheckResult{
				Name:    ScriptHashCheckName,
				Version: ScriptHashCheckVersion,
				Error:   nil,
				Abort:   false,
				Mfa:     true,
				Details: fmt.Sprintf(ScriptHashCheckMfaScriptDetails+ScriptHashCheckUpdatedScriptDetails,
					"testJob", "sha256:v7gXKQ__R6c76dShfVJoe8gvPOc-TJxHhv3QDdRp5kw=", "testJob"),
			},
		},
	}

	for testNum, test := range scriptHashCheckTests {
		var wg sync.WaitGroup
		channel := make(chan checks.CheckResult, 1)

		wg.Add(1)
		go test.scriptHashCheck.Check(channel, &wg, test.whitelist)

		wg.Wait()
		close(channel)

		result := <-channel

		if !result.CompareCheckResult(test.expectedResult) {
			t.Errorf("\n%d) Results not equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedResult, result)
		}
	}
}

func TestHashScript(t *testing.T) {
	hashScriptTests := []struct {
		name            string
		scriptHashCheck ScriptHashCheck
		expectedHash    string
	}{
		{
			name: "hash single line script",
			scriptHashCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   true,
				scriptLines: []string{
					`$ echo "this script is for testing"`,
					`$ curl -X POST -H yourdomain.com"`,
				},
			},
			expectedHash: "sha256:v7gXKQ__R6c76dShfVJoe8gvPOc-TJxHhv3QDdRp5kw=",
		},
		{
			name: "hash multi line script",
			scriptHashCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   true,
				scriptLines: []string{
					`$ echo "test script"`,
				},
			},
			expectedHash: "sha256:D3Xo55v3W4MNn6U3hSH5S6ShRPryxaSKOnE0tQk-O00=",
		},
		{
			name: "hash empty script",
			scriptHashCheck: ScriptHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   true,
				scriptLines: []string{},
			},
			expectedHash: "",
		},
	}

	for testNum, test := range hashScriptTests {
		hash := test.scriptHashCheck.hashScript()
		if hash != test.expectedHash {
			t.Errorf("\n%d) Resulting Hashes not equal -\nexpected: (%s)\ngot: (%s)", testNum, test.expectedHash, hash)
		}
	}
}

func TestIsValidForCheckType(t *testing.T) {
	isValidForCheckTypeTests := []struct {
		name      string
		checkType uint
		valid     bool
	}{
		{
			name:      "script hash check should be valid for all checks",
			checkType: checks.All,
			valid:     true,
		},
		{
			name:      "script hash check should be valid for script checks",
			checkType: checks.ScriptCheck,
			valid:     true,
		},
		{
			name:      "script hash check should not be valid for image checks",
			checkType: checks.ImageCheck,
			valid:     false,
		},
	}

	scriptHashCheck := ScriptHashCheck{}
	for testNum, test := range isValidForCheckTypeTests {
		valid := scriptHashCheck.IsValidForCheckType(test.checkType)
		if valid != test.valid {
			t.Errorf("\n%d) unexpected valid for check type -\nexpected: (%t)\ngot: (%t)", testNum, test.valid, valid)
		}
	}
}

func TestIsValidForPlatform(t *testing.T) {
	isValidForPlatformTests := []struct {
		name     string
		platform string
		valid    bool
	}{
		{
			name:     "script hash check should be valid for github",
			platform: "github",
			valid:    true,
		},
		{
			name:     "script hash check should be valid for gitlab",
			platform: "gitlab",
			valid:    true,
		},
	}

	scriptHashCheck := ScriptHashCheck{}
	for testNum, test := range isValidForPlatformTests {
		valid := scriptHashCheck.IsValidForPlatform(test.platform)
		if valid != test.valid {
			t.Errorf("\n%d) unexpected valid for platform -\nexpected: (%t)\ngot: (%t)", testNum, test.valid, valid)
		}
	}
}
