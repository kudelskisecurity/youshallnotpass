package mfarequired

import (
	"reflect"
	"sync"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

func TestNewMfaRequiredCheck(t *testing.T) {
	newMfaRequiredCheckTests := []struct {
		name          string
		config        config.CheckConfig
		expectedCheck MfaRequiredCheck
	}{
		{
			name: "mfa required check no check type specified",
			config: config.CheckConfig{
				Name:    "MfaRequired",
				Options: map[string]interface{}{},
			},
			expectedCheck: MfaRequiredCheck{
				jobName:        "testJob",
				validCheckType: checks.All,
			},
		},
		{
			name: "mfa required check image check specified",
			config: config.CheckConfig{
				Name: "MfaRequired",
				Options: map[string]interface{}{
					"checkType": "image",
				},
			},
			expectedCheck: MfaRequiredCheck{
				jobName:        "testJob",
				validCheckType: checks.ImageCheck,
			},
		},
		{
			name: "mfa required check script check specified",
			config: config.CheckConfig{
				Name: "MfaRequired",
				Options: map[string]interface{}{
					"checkType": "script",
				},
			},
			expectedCheck: MfaRequiredCheck{
				jobName:        "testJob",
				validCheckType: checks.ScriptCheck,
			},
		},
	}

	jobName := "testJob"
	for testNum, test := range newMfaRequiredCheckTests {
		check := NewMfaRequiredCheck(test.config, jobName)

		if !reflect.DeepEqual(check, test.expectedCheck) {
			t.Errorf("\n%d) checks not equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedCheck, check)
		}
	}
}
func TestMfaRequiredCheck(t *testing.T) {
	mfaRequiredCheckTests := []struct {
		name             string
		mfaRequiredCheck MfaRequiredCheck
		expectedResult   checks.CheckResult
	}{
		{
			name: "perform mfa required check",
			mfaRequiredCheck: MfaRequiredCheck{
				jobName:        "testJob",
				validCheckType: checks.All,
			},
			expectedResult: checks.CheckResult{
				Name:    MfaRequiredCheckName,
				Version: MfaRequiredCheckVersion,
				Error:   nil,
				Mfa:     true,
				Details: MfaRequiredCheckDetails,
			},
		},
	}

	whitelist := whitelist.Whitelist{}
	for testNum, test := range mfaRequiredCheckTests {
		var wg sync.WaitGroup
		channel := make(chan checks.CheckResult, 1)

		wg.Add(1)
		go test.mfaRequiredCheck.Check(channel, &wg, whitelist)

		wg.Wait()
		close(channel)

		result := <-channel

		if !result.CompareCheckResult(test.expectedResult) {
			t.Errorf("\n%d) Results not equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedResult, result)
		}
	}
}

func TestIsValidForCheckType(t *testing.T) {
	isValidForCheckTypeTests := []struct {
		name             string
		mfaRequiredCheck MfaRequiredCheck
		checkType        uint
		valid            bool
	}{
		{
			name: "mfa required checks valid for all checks should be valid for all checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.All,
			},
			checkType: checks.All,
			valid:     true,
		},
		{
			name: "mfa required checks valid for all checks should be valid for image checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.All,
			},
			checkType: checks.ImageCheck,
			valid:     true,
		},
		{
			name: "mfa required checks valid for all checks should be valid for script checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.All,
			},
			checkType: checks.ScriptCheck,
			valid:     true,
		},
		{
			name: "mfa required checks valid for image checks should be valid for all checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.ImageCheck,
			},
			checkType: checks.All,
			valid:     true,
		},
		{
			name: "mfa required checks valid for image checks should be valid for image checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.ImageCheck,
			},
			checkType: checks.ImageCheck,
			valid:     true,
		},
		{
			name: "mfa required checks valid for image checks should not be valid for script checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.ImageCheck,
			},
			checkType: checks.ScriptCheck,
			valid:     false,
		},
		{
			name: "mfa required checks valid for script checks should be valid for all checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.ScriptCheck,
			},
			checkType: checks.All,
			valid:     true,
		},
		{
			name: "mfa required checks valid for script checks should be valid for script checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.ScriptCheck,
			},
			checkType: checks.ScriptCheck,
			valid:     true,
		},
		{
			name: "mfa required checks valid for script checks should not be valid for image checks",
			mfaRequiredCheck: MfaRequiredCheck{
				validCheckType: checks.ScriptCheck,
			},
			checkType: checks.ImageCheck,
			valid:     false,
		},
	}

	for testNum, test := range isValidForCheckTypeTests {
		valid := test.mfaRequiredCheck.IsValidForCheckType(test.checkType)
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
			name:     "mfa required check should be valid for github",
			platform: "github",
			valid:    true,
		},
		{
			name:     "mfa required check should be valid for gitlab",
			platform: "gitlab",
			valid:    true,
		},
	}

	mfaRequiredCheck := MfaRequiredCheck{}
	for testNum, test := range isValidForPlatformTests {
		valid := mfaRequiredCheck.IsValidForPlatform(test.platform)
		if valid != test.valid {
			t.Errorf("\n%d) unexpected valid for platform -\nexpected: (%t)\ngot: (%t)", testNum, test.valid, valid)
		}
	}
}
