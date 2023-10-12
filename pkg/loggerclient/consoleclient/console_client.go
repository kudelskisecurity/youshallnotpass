package consoleclient

import (
	"fmt"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
)

var (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

type ConsoleClient struct {
	ciJobName         string
	vaultExternalAddr string
	vaultRoot         string
	ciProjectPath     string
}

func NewConsoleClient(ciJobName string, vaultExternalAddr string, vaultRoot string, ciProjectPath string) *ConsoleClient {
	return &ConsoleClient{
		ciJobName,
		vaultExternalAddr,
		vaultRoot,
		ciProjectPath,
	}
}

func (c *ConsoleClient) LogRecoverableError(err error) {
	fmt.Print(err.Error())
}

func (c *ConsoleClient) SendMFAInstructions(message string) error {
	print(message)
	return nil
}

func (c *ConsoleClient) LogCheckResults(results []checks.CheckResult) error {
	message := "---------------------------------------------------------------------\n"
	message += fmt.Sprintf("%s/ui/vault/secrets/%s - %s\n", c.vaultExternalAddr, c.vaultRoot+"/show/"+c.ciProjectPath+"/whitelist", c.ciJobName)
	message += "---------------------------------------------------------------------\n"
	message += "   Name   |   Version   |   Error   |   Abort   |   Mfa   |   Details\n"
	message += "---------------------------------------------------------------------\n"
	for _, result := range results {
		errStr := ""
		if result.Error != nil {
			errStr = result.Error.Error()
		}
		message += fmt.Sprintf("%s | %s | %s | %t | %t | %s\n", result.Name, result.Version,
			errStr, result.Abort, result.Mfa, result.Details)
		message += "---------------------------------------------------------------------\n"
	}
	fmt.Printf("\n\n%s\n\n", message)
	return nil
}

func (c *ConsoleClient) LogPrevalidated() error {
	fmt.Printf("  ✅ CI/CD for %s has been prevalidated", c.ciJobName)
	return nil
}

func (c *ConsoleClient) LogSuccessfulExecution() error {
	fmt.Printf("  ✅ %sSuccessful YouShallNotPass Check for Job: %s%s\n", string(colorGreen), c.ciJobName, string(colorReset))
	return nil
}

func (c *ConsoleClient) LogFailedExecution(reason string) error {
	fmt.Printf("  ❌ %sUnsuccessful YouShallNotPass Check for Job: %s%s\n", string(colorRed), c.ciJobName, string(colorReset))
	return fmt.Errorf(reason)
}
