package config

import "encoding/json"

type LoggerConfig struct {
	Name    string                 `json:"name,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type NamespaceConfig struct {
	LoggerConfig LoggerConfig `json:"logger,omitempty"`
}

type CheckConfig struct {
	Name    string                 `json:"name,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type JobConfig struct {
	JobName string        `json:"jobName,omitempty"`
	Checks  []CheckConfig `json:"checks,omitempty"`
}

type ProjectConfig struct {
	Jobs []JobConfig `json:"jobs,omitempty"`
}

var DefaultNamespaceConfig = NamespaceConfig{
	LoggerConfig: LoggerConfig{
		Name: "console",
	},
}

var DefaultProjectConfig = ProjectConfig{
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
				{
					Name: "scriptHash",
					Options: map[string]interface{}{
						"mfaOnFail": true,
					},
				},
			},
		},
	},
}

func ParseNamespaceConfig(jsonText []byte) (NamespaceConfig, error) {
	if len(jsonText) == 0 {
		return DefaultNamespaceConfig, nil
	}

	var namespaceConfig NamespaceConfig

	err := json.Unmarshal([]byte(jsonText), &namespaceConfig)
	if err != nil {
		return DefaultNamespaceConfig, err
	}

	return namespaceConfig, nil
}

func ParseProjectConfig(jsonText []byte) (ProjectConfig, error) {
	if len(jsonText) == 0 {
		return DefaultProjectConfig, nil
	}

	var projectConfig ProjectConfig

	err := json.Unmarshal([]byte(jsonText), &projectConfig)
	if err != nil {
		return DefaultProjectConfig, err
	}

	return projectConfig, nil
}

func (c *ProjectConfig) GetJobConfig(jobName string) []CheckConfig {
	if len(jobName) == 0 {
		jobName = "default"
	}

	for _, job := range c.Jobs {
		if job.JobName == jobName {
			return job.Checks
		}
	}

	for _, job := range c.Jobs {
		if job.JobName == "default" {
			return job.Checks
		}
	}

	return DefaultProjectConfig.Jobs[0].Checks
}
