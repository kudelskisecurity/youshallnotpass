package config

import (
	"reflect"
	"testing"
)

func TestParseNamespaceConfig(t *testing.T) {
	parseNamespaceConfigTests := []struct {
		name           string
		jsonText       string
		expectedConfig NamespaceConfig
		errorExpected  bool
	}{
		{
			name: "parse mattermost client configuration",
			jsonText: `
{
	"logger": {
		"name": "mattermost",
		"options": {
			"url": "http://127.0.0.1:8000/mattermost",
			"token": "1234567890",
			"channelId": "1234567890"
		}
	}
}`,
			expectedConfig: NamespaceConfig{
				LoggerConfig{
					Name: "mattermost",
					Options: map[string]interface{}{
						"url":       "http://127.0.0.1:8000/mattermost",
						"token":     "1234567890",
						"channelId": "1234567890",
					},
				},
			},
			errorExpected: false,
		},
		{
			name: "parse console client configuration",
			jsonText: `
{
	"logger": {
		"name": "console"
	}
}`,
			expectedConfig: NamespaceConfig{
				LoggerConfig{
					Name: "console",
				},
			},
			errorExpected: false,
		},
		{
			name:           "parse default namespace configuration",
			jsonText:       ``,
			expectedConfig: DefaultNamespaceConfig,
			errorExpected:  false,
		},
		{
			name: "parse invalid json namespace configuration",
			jsonText: `
{
	"logger",: { "tests" },
}`,
			expectedConfig: DefaultNamespaceConfig,
			errorExpected:  true,
		},
	}

	for testNum, test := range parseNamespaceConfigTests {
		namespaceConfig, err := ParseNamespaceConfig([]byte(test.jsonText))
		if !test.errorExpected && err != nil {
			t.Errorf("\n%d) unexpected error: %s", testNum, err.Error())
		} else if test.errorExpected && err == nil {
			t.Errorf("\n%d) expected an error", testNum)
		}

		if !reflect.DeepEqual(namespaceConfig, test.expectedConfig) {
			t.Errorf("\n%d) Namespace Configurations not equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedConfig, namespaceConfig)
		}
	}
}

func TestParseProjectConfig(t *testing.T) {
	praseProjectConfigTests := []struct {
		name           string
		jsonText       string
		expectedConfig ProjectConfig
		errorExpected  bool
	}{
		{
			name: "parse single job config",
			jsonText: `
{
	"jobs": [
		{
			"jobName": "job1",
			"checks": [
				{
					"name": "imageHash",
					"options": {
						"abortOnFail": true
					}
				}
			]
		}
	]
}`,
			expectedConfig: ProjectConfig{
				Jobs: []JobConfig{
					{
						JobName: "job1",
						Checks: []CheckConfig{
							{
								Name: "imageHash",
								Options: map[string]interface{}{
									"abortOnFail": true,
								},
							},
						},
					},
				},
			},
			errorExpected: false,
		},
		{
			name:           "parse no job configs",
			jsonText:       "",
			expectedConfig: DefaultProjectConfig,
			errorExpected:  false,
		},
		{
			name: "parse many job configs",
			jsonText: `
{
	"jobs": [
		{
			"jobName": "testJobOne",
			"checks": [
				{
					"name": "imageHash",
					"options": {
						"abortOnFail": true
					}
				}
			]
		},
		{
			"jobName": "testJobTwo",
			"checks": [
				{
					"name": "scriptHash",
					"options": {
						"mfaOnFail": true
					}
				}
			]
		}
	]
}`,
			expectedConfig: ProjectConfig{
				Jobs: []JobConfig{
					{
						JobName: "testJobOne",
						Checks: []CheckConfig{
							{
								Name: "imageHash",
								Options: map[string]interface{}{
									"abortOnFail": true,
								},
							},
						},
					},
					{
						JobName: "testJobTwo",
						Checks: []CheckConfig{
							{
								Name: "scriptHash",
								Options: map[string]interface{}{
									"mfaOnFail": true,
								},
							},
						},
					},
				},
			},
			errorExpected: false,
		},
		{
			name: "parse invalid json",
			jsonText: `
{
	"test",,:{should"fail}
}`,
			expectedConfig: DefaultProjectConfig,
			errorExpected:  true,
		},
	}

	for testNum, test := range praseProjectConfigTests {
		config, err := ParseProjectConfig([]byte(test.jsonText))
		if !test.errorExpected && err != nil {
			t.Errorf("\n%d) unexpected error: %s", testNum, err.Error())
		} else if test.errorExpected && err == nil {
			t.Errorf("\n%d) expected an error", testNum)
		}

		if !reflect.DeepEqual(config, test.expectedConfig) {
			t.Errorf("\n%d) Project Configurations not equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedConfig, config)
		}
	}
}

func TestGetJobConfig(t *testing.T) {
	getJobConfigTests := []struct {
		name           string
		jobName        string
		projectConfig  ProjectConfig
		expectedChecks []CheckConfig
	}{
		{
			name:    "test get job checks",
			jobName: "testJob",
			projectConfig: ProjectConfig{
				Jobs: []JobConfig{
					{
						JobName: "testJob",
						Checks: []CheckConfig{
							{
								Name: "imageHash",
								Options: map[string]interface{}{
									"abortOnFail": true,
								},
							},
						},
					},
				},
			},
			expectedChecks: []CheckConfig{
				{
					Name: "imageHash",
					Options: map[string]interface{}{
						"abortOnFail": true,
					},
				},
			},
		},
		{
			name:    "test get default job checks",
			jobName: "testJob",
			projectConfig: ProjectConfig{
				Jobs: []JobConfig{
					{
						JobName: "default",
						Checks: []CheckConfig{
							{
								Name: "imageHash",
								Options: map[string]interface{}{
									"abortOnFail": true,
								},
							},
						},
					},
				},
			},
			expectedChecks: []CheckConfig{
				{
					Name: "imageHash",
					Options: map[string]interface{}{
						"abortOnFail": true,
					},
				},
			},
		},
		{
			name:           "test get youshallnotpass default job checks",
			jobName:        "testJob",
			projectConfig:  ProjectConfig{},
			expectedChecks: DefaultProjectConfig.Jobs[0].Checks,
		},
		{
			name:           "test get empty job name",
			jobName:        "",
			projectConfig:  ProjectConfig{},
			expectedChecks: DefaultProjectConfig.Jobs[0].Checks,
		},
	}

	for testNum, test := range getJobConfigTests {
		checks := test.projectConfig.GetJobConfig(test.jobName)
		if !reflect.DeepEqual(test.expectedChecks, checks) {
			t.Errorf("\n%d) Checks are not equal -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedChecks, checks)
		}
	}
}
