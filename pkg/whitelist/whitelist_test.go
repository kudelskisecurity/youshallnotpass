package whitelist

import (
	"reflect"
	"testing"
)

func TestContainsImage(t *testing.T) {
	containsImageTests := []struct {
		name      string
		whitelist Whitelist
		image     string
		foundImg  bool
		err       error
	}{
		{
			name: "test contains image (image in whitelist)",
			whitelist: Whitelist{
				AllowedImages: []string{
					"alpine:3.12.7@sha256:a9c28c813336ece5bb98b36af5b66209ed777a394f4f856c6e62267790883820",
				},
			},
			image:    "alpine:3.12.7@sha256:a9c28c813336ece5bb98b36af5b66209ed777a394f4f856c6e62267790883820",
			foundImg: true,
			err:      nil,
		},
		{
			name: "test contains image (image not in whitelist)",
			whitelist: Whitelist{
				AllowedImages: []string{
					"alpine:3.12.7@sha256:a9c28c813336ece5bb98b36af5b66209ed777a394f4f856c6e62267790883820",
				},
			},
			image:    "alpine:3.12.7@sha256:a9c28c813336ece5bb98b36af5b66209ed777a394f4f856c6e62267790883abc",
			foundImg: false,
			err:      nil,
		},
		{
			name: "test contains image invalid image (error)",
			whitelist: Whitelist{
				AllowedImages: []string{
					"alpine:3.12.7@sha256:a9c28c813336ece5bb98b36af5b66209ed777a394f4f856c6e62267790883820",
				},
			},
			image:    "alping:3.12.7@sha256:abc",
			foundImg: false,
			err:      ErrShaNotValidSha,
		},
	}

	for testNum, test := range containsImageTests {
		got, err := test.whitelist.ContainsImage(test.image)
		if err != test.err {
			t.Errorf("\n%d) unexpected error -\nexpected: (%+v)\ngot: (%+v)", testNum, test.err, err)
		}

		if got != test.foundImg {
			t.Errorf("\n%d) unexpected found image -\nexpected: (%t)\ngot: (%t)", testNum, test.foundImg, got)
		}
	}
}

func TestContainsScript(t *testing.T) {
	containsScriptTests := []struct {
		name        string
		whitelist   Whitelist
		scriptSha   string
		foundScript bool
	}{
		{
			name: "test contains script (script in whitelist)",
			whitelist: Whitelist{
				AllowedScripts: []string{
					"automatic_job@sha256:oNey8xJbYXyuxWr7Wla8tMexCTy7s82k6U1uwp4tFEY=",
				},
			},
			scriptSha:   "sha256:oNey8xJbYXyuxWr7Wla8tMexCTy7s82k6U1uwp4tFEY=",
			foundScript: true,
		},
		{
			name: "test contains script (script not in whitelist)",
			whitelist: Whitelist{
				AllowedScripts: []string{
					"automatic_job@sha256:oNey8xJbYXyuxWr7Wla8tMexCTy7s82k6U1uwp4tFEY=",
				},
			},
			scriptSha:   "build_job@sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
			foundScript: false,
		},
		{
			name: "test contains script (invalid script hash format)",
			whitelist: Whitelist{
				AllowedScripts: []string{
					"automatic_job@sha256:oNey8xJbY=",
				},
			},
			scriptSha:   "sha256:oNey8xJbYXyuxWr7Wla8tMexCTy7s82k6U1uwp4tFEY=",
			foundScript: false,
		},
	}

	for testNum, test := range containsScriptTests {
		got := test.whitelist.ContainsScript(test.scriptSha)
		if got != test.foundScript {
			t.Errorf("\n%d) unexpected found script -\nexpected: (%t)\ngot: (%t)", testNum, test.foundScript, got)
		}
	}
}

func TestContainsJobName(t *testing.T) {
	containsJobNameTests := []struct {
		name         string
		whitelist    Whitelist
		jobName      string
		foundJobName bool
		foundJobHash string
	}{
		{
			name: "whitelist contains job name",
			whitelist: Whitelist{
				AllowedScripts: []string{
					"testJob@sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
				},
			},
			jobName:      "testJob",
			foundJobName: true,
			foundJobHash: "sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
		},
		{
			name: "whitelist does not contain job name",
			whitelist: Whitelist{
				AllowedScripts: []string{
					"build_job@sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
				},
			},
			jobName:      "testJob",
			foundJobName: false,
			foundJobHash: "",
		},
	}

	for testNum, test := range containsJobNameTests {
		got, hash := test.whitelist.ContainsJobName(test.jobName)
		if got != test.foundJobName {
			t.Errorf("\n%d) unexpected found job name -\nexpected: (%t)\ngot: (%t)", testNum, test.foundJobName, got)
		}

		if hash != test.foundJobHash {
			t.Errorf("\n%d) unexpected hash found -\nexpected: (%s)\ngot: (%s)", testNum, test.foundJobHash, hash)
		}
	}
}

func TestAddWhitelist(t *testing.T) {
	addWhitelistTests := []struct {
		name              string
		whitelistOne      Whitelist
		whitelistTwo      Whitelist
		expectedWhitelist Whitelist
	}{
		{
			name:              "add empty whitelists",
			whitelistOne:      Whitelist{},
			whitelistTwo:      Whitelist{},
			expectedWhitelist: Whitelist{},
		},
		{
			name: "add whitelists allowed images and scripts",
			whitelistOne: Whitelist{
				AllowedImages: []string{
					"imageOne",
					"imageTwo",
				},
				AllowedScripts: []string{
					"scriptOne",
					"scriptTwo",
				},
			},
			whitelistTwo: Whitelist{
				AllowedImages: []string{
					"imageThree",
				},
				AllowedScripts: []string{
					"scriptThree",
				},
			},
			expectedWhitelist: Whitelist{
				AllowedImages: []string{
					"imageOne",
					"imageTwo",
					"imageThree",
				},
				AllowedScripts: []string{
					"scriptOne",
					"scriptTwo",
					"scriptThree",
				},
			},
		},
	}

	for testNum, test := range addWhitelistTests {
		test.whitelistOne.AddWhitelist(test.whitelistTwo)
		if !reflect.DeepEqual(test.whitelistOne, test.expectedWhitelist) {
			t.Errorf("\n%d) expected whitelists to be equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.whitelistOne, test.expectedWhitelist)
		}
	}
}

func TestGetJobName(t *testing.T) {
	getJobNameTests := []struct {
		name            string
		whitelistScript string
		jobName         string
		err             error
	}{
		{
			name:            "get job name normal",
			whitelistScript: "build-script@sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
			jobName:         "build-script",
			err:             nil,
		},
		{
			name:            "get job name invalid name (i.e. empty name)",
			whitelistScript: "@sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
			jobName:         "",
			err:             ErrInvalidWhitelistJob,
		},
	}

	for testNum, test := range getJobNameTests {
		jobName, err := getJobName(test.whitelistScript)
		if err != test.err {
			t.Errorf("\n%d) unexpected error -\nexpected: (%+v)\ngot: (%+v)", testNum, test.err, err)
		}

		if jobName != test.jobName {
			t.Errorf("\n%d) unexpected job name -\nexpected: (%s)\ngot: (%s)", testNum, test.jobName, jobName)
		}
	}
}

func TestGetScriptSha(t *testing.T) {
	getScriptShaTests := []struct {
		name            string
		whitelistScript string
		sha             string
		err             error
	}{
		{
			name:            "get script sha expected whitelist script",
			whitelistScript: "build-script@sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
			sha:             "sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
			err:             nil,
		},
		{
			name:            "get script sha no sha",
			whitelistScript: "build-script",
			sha:             "",
			err:             ErrShaNotPresent,
		},
		{
			name:            "get script sha no job name",
			whitelistScript: "sha256:1NmXdCi0PhRNMKX91bypxJKJhzB1nUQqtRowPbgpqqE=",
			sha:             "",
			err:             ErrShaNotPresent,
		},
		{
			name:            "get script sha empty whitelist script",
			whitelistScript: "",
			sha:             "",
			err:             ErrShaNotPresent,
		},
		{
			name:            "get script sha invalid sha",
			whitelistScript: "build-script@sha256:1NmXdCi0PhRNMK",
			sha:             "",
			err:             ErrShaNotValidSha,
		},
	}

	for testNum, test := range getScriptShaTests {
		sha, err := getScriptSha(test.whitelistScript)
		if err != test.err {
			t.Errorf("\n%d) unexpected error -\nexpected: (%+v)\ngot: (%+v)", testNum, test.err, err)
		}

		if sha != test.sha {
			t.Errorf("\n%d) unexpected sha found -\nexpected: (%s)\ngot: (%s)", testNum, test.sha, sha)
		}
	}
}

func TestGetImageSha(t *testing.T) {
	getImageShaTests := []struct {
		name           string
		whitelistImage string
		sha            string
		err            error
	}{
		{
			name:           "get image sha expected whitelist",
			whitelistImage: "alpine:latest@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
			sha:            "sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
			err:            nil,
		},
		{
			name:           "get image sha no sha",
			whitelistImage: "alpine:3.13.1",
			sha:            "",
			err:            ErrShaNotPresent,
		},
		{
			name:           "get image sha not image name",
			whitelistImage: "sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748",
			sha:            "",
			err:            ErrShaNotPresent,
		},
		{
			name:           "get image sha empty whitelist image",
			whitelistImage: "",
			sha:            "",
			err:            ErrShaNotPresent,
		},
		{
			name:           "get image sha invalid sha",
			whitelistImage: "alpine:latest@sha256:def822f9851ca422481ec6f",
			sha:            "",
			err:            ErrShaNotValidSha,
		},
	}

	for testNum, test := range getImageShaTests {
		sha, err := getImageSha(test.whitelistImage)
		if err != test.err {
			t.Errorf("\n%d) unexpected error -\nexpected: (%+v)\ngot: (%+v)", testNum, test.err, err)
		}

		if sha != test.sha {
			t.Errorf("\n%d) unexpected sha -\nexpected: (%s)\ngot: (%s)", testNum, test.sha, sha)
		}
	}
}
