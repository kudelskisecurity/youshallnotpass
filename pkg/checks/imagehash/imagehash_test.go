package imagehash

import (
	"reflect"
	"sync"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

func TestNewImageHashCheck(t *testing.T) {
	newImageHashCheckTests := []struct {
		name          string
		config        config.CheckConfig
		expectedCheck ImageHashCheck
	}{
		{
			name: "test abortOnFail true and mfaOnFail true",
			config: config.CheckConfig{
				Name: "imageHash",
				Options: map[string]interface{}{
					"abortOnFail": true,
					"mfaOnFail":   true,
				},
			},
			expectedCheck: ImageHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   true,
				image:       "",
			},
		},
		{
			name: "test abortOnFail false and mfaOnFail false",
			config: config.CheckConfig{
				Name: "imageHash",
				Options: map[string]interface{}{
					"abortOnFail": false,
					"mfaOnFail":   false,
				},
			},
			expectedCheck: ImageHashCheck{
				jobName:     "testJob",
				abortOnFail: false,
				mfaOnFail:   false,
				image:       "",
			},
		},
		{
			name: "test abortOnFail missing and mfaOnFail missing",
			config: config.CheckConfig{
				Name:    "imageHash",
				Options: map[string]interface{}{},
			},
			expectedCheck: ImageHashCheck{
				jobName:     "testJob",
				abortOnFail: false,
				mfaOnFail:   false,
				image:       "",
			},
		},
	}

	jobName := "testJob"
	image := ""
	for _, test := range newImageHashCheckTests {
		imageHashCheck := NewImageHashCheck(test.config, jobName, image)
		if !reflect.DeepEqual(imageHashCheck, test.expectedCheck) {
			t.Errorf("checks not equal - expected: (%+v) - got: (%+v)", test.expectedCheck, imageHashCheck)
		}
	}
}
func TestImageHashCheck(t *testing.T) {
	imageHashCheckTests := []struct {
		name           string
		imageHashCheck ImageHashCheck
		whitelist      whitelist.Whitelist
		expectedResult checks.CheckResult
	}{
		{
			name: "test image in whitelist - success",
			imageHashCheck: ImageHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   true,
				image:       "alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
			},
			whitelist: whitelist.Whitelist{
				AllowedImages: []string{
					"alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
				},
			},
			expectedResult: checks.CheckResult{
				Name:    "Image Hash Check",
				Version: "1.0.0",
				Error:   nil,
				Abort:   false,
				Mfa:     false,
				Details: ImageHashCheckSuccess,
			},
		},
		{
			name: "test image not in whitelist - mfaOnFail",
			imageHashCheck: ImageHashCheck{
				jobName:     "testJob",
				abortOnFail: false,
				mfaOnFail:   true,
				image:       "alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
			},
			whitelist: whitelist.Whitelist{
				AllowedImages: []string{
					"alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9fxyz",
				},
			},
			expectedResult: checks.CheckResult{
				Name:    "Image Hash Check",
				Version: "1.0.0",
				Error:   nil,
				Abort:   false,
				Mfa:     true,
				Details: ImageHashCheckFailedMfaRequired,
			},
		},
		{
			name: "test image not in whitelist - abortOnfail",
			imageHashCheck: ImageHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   false,
				image:       "alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
			},
			whitelist: whitelist.Whitelist{
				AllowedImages: []string{
					"alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9fxyz",
				},
			},
			expectedResult: checks.CheckResult{
				Name:    "Image Hash Check",
				Version: "1.0.0",
				Error:   nil,
				Abort:   true,
				Mfa:     false,
				Details: ImageHashCheckFailedAbort,
			},
		},
		{
			name: "test image hash check no image",
			imageHashCheck: ImageHashCheck{
				jobName:     "testJob",
				abortOnFail: true,
				mfaOnFail:   true,
				image:       "",
			},
			whitelist: whitelist.Whitelist{
				AllowedImages: []string{
					"alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
				},
			},
			expectedResult: checks.CheckResult{
				Name:    "Image Hash Check",
				Version: "1.0.0",
				Error:   whitelist.ErrShaNotPresent,
				Abort:   true,
				Mfa:     false,
				Details: ImageHashCheckError,
			},
		},
		{
			name: "test image hash check invalid hash",
			imageHashCheck: ImageHashCheck{
				jobName:     "testJob",
				abortOnFail: false,
				mfaOnFail:   true,
				image:       "testImage@sha256:XYZ",
			},
			whitelist: whitelist.Whitelist{
				AllowedImages: []string{
					"alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
				},
			},
			expectedResult: checks.CheckResult{
				Name:    "Image Hash Check",
				Version: "1.0.0",
				Error:   whitelist.ErrShaNotValidSha,
				Abort:   false,
				Mfa:     true,
				Details: ImageHashCheckError,
			},
		},
	}

	for testNum, test := range imageHashCheckTests {
		var wg sync.WaitGroup
		channel := make(chan checks.CheckResult, 1)

		wg.Add(1)
		go test.imageHashCheck.Check(channel, &wg, test.whitelist)

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
		name      string
		checkType uint
		valid     bool
	}{
		{
			name:      "image hash check should be valid for all checks",
			checkType: checks.All,
			valid:     true,
		},
		{
			name:      "image hash check should be valid for image checks",
			checkType: checks.ImageCheck,
			valid:     true,
		},
		{
			name:      "image hash check should not be valid for script checks",
			checkType: checks.ScriptCheck,
			valid:     false,
		},
	}

	imageHashCheck := ImageHashCheck{}
	for testNum, test := range isValidForCheckTypeTests {
		valid := imageHashCheck.IsValidForCheckType(test.checkType)
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
			name:     "image hash check should be valid for github",
			platform: "github",
			valid:    true,
		},
		{
			name:     "image hash check should be valid for gitlab",
			platform: "gitlab",
			valid:    true,
		},
	}

	imageHashCheck := ImageHashCheck{}
	for testNum, test := range isValidForPlatformTests {
		valid := imageHashCheck.IsValidForPlatform(test.platform)
		if valid != test.valid {
			t.Errorf("\n%d) unexpected valid for platform -\nexpected: (%t)\ngot: (%t)", testNum, test.valid, valid)
		}
	}
}
