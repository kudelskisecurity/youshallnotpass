package loggerclient

import (
	"fmt"
	"strings"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/loggerclient/consoleclient"
	"github.com/kudelskisecurity/youshallnotpass/pkg/loggerclient/mattermostclient"
)

type LoggerClient interface {
	LogRecoverableError(error)
	SendMFAInstructions(string) error
	LogCheckResults([]checks.CheckResult) error
	LogPrevalidated() error
	LogSuccessfulExecution() error
	LogFailedExecution(string) error
}

func ParseNotifyClient(c config.LoggerConfig, ciJobName string, vaultExternalAddr string, vaultRoot string, ciProjectPath string) (LoggerClient, error) {
	if strings.ToLower(c.Name) == "mattermost" {
		url, exists := c.Options["url"].(string)
		if !exists {
			return nil, fmt.Errorf("to use mattermost as a logger please include an instance url int vault")
		}

		token, exists := c.Options["token"].(string)
		if !exists {
			return nil, fmt.Errorf("to use mattermost as a logger please include a token in vault")
		}

		channelId, exists := c.Options["channelId"].(string)
		if !exists {
			return nil, fmt.Errorf("to use mattermost as a logger please include a channel id in vault")
		}

		return mattermostclient.NewMattermostClient(url, token, channelId, ciJobName, vaultExternalAddr, vaultRoot, ciProjectPath)
	}

	return consoleclient.NewConsoleClient(ciJobName, vaultExternalAddr, vaultRoot, ciProjectPath), nil
}
