package hashicorpclient

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/urfave/cli/v2"
)

// This is more integration testing than unit testing because the hashicorp vault
// client is pretty much known to work we're just wrapping it

func TestInitVaultClientValidVaultToken(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "http://0.0.0.0:8200", "vault address")
	contextSet.String("vault-token", "1234567890", "vault token")
	contextSet.String("jwt-token", "", "jwt token")
	contextSet.String("vault-role", "youshallnotpass-demo", "vault role")
	contextSet.String("vault-root", "cicd", "vault root")
	contextSet.String("ci-project-path", "youshallnotpass/demo", "ci project path")
	context := cli.NewContext(nil, contextSet, nil)

	_, err := InitVaultClient(context)
	if err != nil {
		t.Errorf("Expected no error, but got %s", err.Error())
		return
	}
}

func TestInitVaultClientNoTokenOrJWT(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "", "vault address")
	contextSet.String("vault-token", "", "authentication token to access the vault")
	contextSet.String("jwt-token", "", "jwt token to authenticat to the vault")
	contextSet.String("vault-role", "", "vault access role")
	context := cli.NewContext(nil, contextSet, nil)

	_, err := InitVaultClient(context)
	if !strings.Contains(err.Error(), "either vaultToken or jwtToken is required") {
		t.Error("Unexpected Error with invalid jwt-token and jwt-token")
		return
	}
}

func TestVaultClientWriteScratch(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "http://0.0.0.0:8200", "vault address")
	contextSet.String("vault-token", "1234567890", "vault token")
	contextSet.String("jwt-token", "", "jwt token")
	contextSet.String("vault-role", "youshallnotpass-demo", "vault role")
	contextSet.String("vault-root", "cicd", "vault root")
	contextSet.String("ci-project-path", "youshallnotpass/demo", "ci project path")
	context := cli.NewContext(nil, contextSet, nil)

	s, err := InitVaultClient(context)
	if err != nil {
		t.Errorf("Expected no error when initializing vault client, but got %s", err.Error())
		return
	}

	path, err := s.WriteScratch(0, "test.user@example.com")
	if err != nil {
		t.Errorf("expected no error when writing scratch, but got %s", err.Error())
		return
	}

	if !strings.Contains(path, "test.user@example.com/") {
		t.Errorf("Unexpected secret path - got %s", path)
		return
	}
}

func TestVaultClientGetNamespaceConfig(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "http://0.0.0.0:8200", "vault address")
	contextSet.String("vault-token", "1234567890", "vault token")
	contextSet.String("jwt-token", "", "jwt token")
	contextSet.String("vault-role", "youshallnotpass-demo", "vault role")
	contextSet.String("vault-root", "cicd", "vault root")
	contextSet.String("ci-project-path", "youshallnotpass/demo", "ci project path")
	context := cli.NewContext(nil, contextSet, nil)

	s, err := InitVaultClient(context)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
		return
	}

	namespaceMount := "cicd/youshallnotpass/youshallnotpass_config"

	namespaceConfig, err := s.GetNamespaceConfig(namespaceMount)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
		return
	}

	expectedConfig := config.NamespaceConfig{
		LoggerConfig: config.LoggerConfig{
			Name: "console",
		},
	}

	if !reflect.DeepEqual(namespaceConfig, expectedConfig) {
		s := fmt.Sprintf("Expected: \n%+v\nGot: \n%+v\n", expectedConfig, namespaceConfig)
		t.Errorf("Unexpected Configuration Obtained: \n%s", s)
		return
	}
}

func TestVaultClientGetProjectConfig(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "http://0.0.0.0:8200", "vault address")
	contextSet.String("vault-token", "1234567890", "vault token")
	contextSet.String("jwt-token", "", "jwt token")
	contextSet.String("vault-role", "youshallnotpass-demo", "vault role")
	contextSet.String("vault-root", "cicd", "vault root")
	contextSet.String("ci-project-path", "youshallnotpass/demo", "ci project path")
	context := cli.NewContext(nil, contextSet, nil)

	s, err := InitVaultClient(context)
	if err != nil {
		t.Errorf("Expected no error when initializing vault client, but got %s", err.Error())
		return
	}

	configMount := "cicd/youshallnotpass/demo/youshallnotpass_config"

	projectConfig, err := s.GetProjectConfig(configMount)
	if err != nil {
		t.Errorf("Expected no error when reading the youshallnotpass config, but got %s", err.Error())
		return
	}

	expectedConfig := config.ProjectConfig{
		Jobs: []config.JobConfig{
			{
				JobName: "user_mfa_job",
				Checks: []config.CheckConfig{
					{
						Name: "mfaRequired",
						Options: map[string]interface{}{
							"checkType": "script",
						},
					},
				},
			}, {
				JobName: "user_mfa_timeout_job",
				Checks: []config.CheckConfig{
					{
						Name: "imageHash",
						Options: map[string]interface{}{
							"abortOnFail": false,
							"mfaOnFail":   true,
						},
					}, {
						Name: "scriptHash",
						Options: map[string]interface{}{
							"abortOnFail": false,
							"mfaOnFail":   true,
						},
					},
				},
			}, {
				JobName: "automatic_job",
				Checks: []config.CheckConfig{
					{
						Name: "imageHash",
						Options: map[string]interface{}{
							"abortOnFail": true,
						},
					}, {
						Name: "scriptHash",
						Options: map[string]interface{}{
							"abortOnFail": true,
						},
					},
				},
			}, {
				JobName: "script_job",
				Checks: []config.CheckConfig{
					{
						Name: "scriptHash",
						Options: map[string]interface{}{
							"abortOnFail": true,
						},
					},
				},
			}, {
				JobName: "fail_job",
				Checks: []config.CheckConfig{
					{
						Name: "imageHash",
						Options: map[string]interface{}{
							"abortOnFail": true,
						},
					}, {
						Name: "scriptHash",
						Options: map[string]interface{}{
							"abortOnFail": true,
						},
					},
				},
			}, {
				JobName: "default",
				Checks: []config.CheckConfig{
					{
						Name: "imageHash",
						Options: map[string]interface{}{
							"abortOnFail": true,
						},
					},
				},
			},
		},
	}

	if len(projectConfig.Jobs) != len(expectedConfig.Jobs) {
		t.Error("different job number")
		return
	}

	for i := range projectConfig.Jobs {
		if projectConfig.Jobs[i].JobName != expectedConfig.Jobs[i].JobName {
			t.Errorf("different job name")
			return
		}

		for j := range projectConfig.Jobs[i].Checks {
			if projectConfig.Jobs[i].Checks[j].Name != expectedConfig.Jobs[i].Checks[j].Name {
				t.Error("different job check name")
				return
			}

			for k := range projectConfig.Jobs[i].Checks[j].Options {
				if projectConfig.Jobs[i].Checks[j].Options[k] != expectedConfig.Jobs[i].Checks[j].Options[k] {
					t.Errorf("different check options")
					return
				}
			}
		}
	}
}

func TestVaultClientReadWhitelist(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "http://0.0.0.0:8200", "vault address")
	contextSet.String("vault-token", "1234567890", "vault token")
	contextSet.String("jwt-token", "", "jwt token")
	contextSet.String("vault-role", "youshallnotpass-demo", "vault role")
	contextSet.String("vault-root", "cicd", "vault root")
	contextSet.String("ci-project-path", "youshallnotpass/demo", "ci project path")
	context := cli.NewContext(nil, contextSet, nil)

	s, err := InitVaultClient(context)
	if err != nil {
		t.Errorf("Expected no error when initializing vault client, but got %s", err.Error())
		return
	}

	whitelistMounts := []string{
		"cicd/youshallnotpass/whitelist",
		"cicd/youshallnotpass/demo/whitelist",
	}

	whitelist, err := s.ReadWhitelists(whitelistMounts[0], whitelistMounts[1])
	if err != nil {
		t.Errorf("Expected no error when reading the whitelist mounts, but got %s", err.Error())
		return
	}

	whitelistOneImage := "alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748"
	contained, _ := whitelist.ContainsImage(whitelistOneImage)
	if !contained {
		t.Errorf("Expected whitelist to contain the image")
		return
	}

	whitelistTwoImage := "alpine:3.12.7@sha256:a9c28c813336ece5bb98b36af5b66209ed777a394f4f856c6e62267790883820"
	whitelistTwoScriptSha := "sha256:Kn9ysqTdXVzh52gp2LNiX5RMNRxdoAQytneeLcNsycQ="
	whitelistTwoScriptShaTwo := "sha256:Ij3eYc5EwfiLD6rPw9qFpN82ydukCduG4bUL9ltQDy4="

	contained, _ = whitelist.ContainsImage(whitelistTwoImage)
	if !contained {
		t.Errorf("Expected whitelist to contain the image")
		return
	}

	if !whitelist.ContainsScript(whitelistTwoScriptSha) {
		t.Error("Expected whitelist to contain script")
		return
	}

	if !whitelist.ContainsScript(whitelistTwoScriptShaTwo) {
		t.Errorf("Expected whitelist to contain script")
		return
	}
}

func TestVaultClientSecretExists(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "http://0.0.0.0:8200", "vault address")
	contextSet.String("vault-token", "1234567890", "vault token")
	contextSet.String("jwt-token", "", "jwt token")
	contextSet.String("vault-role", "youshallnotpass-demo", "vault role")
	contextSet.String("vault-root", "cicd", "vault root")
	contextSet.String("ci-project-path", "youshallnotpass/demo", "ci project path")
	context := cli.NewContext(nil, contextSet, nil)

	secretMountScratch := "cicd/youshallnotpass/demo/scratch"

	s, err := InitVaultClient(context)
	if err != nil {
		t.Errorf("Expected no error when initializing vault client, but got %s", err.Error())
		return
	}

	path, err := s.WriteScratch(0, "test.user@example.com")
	if err != nil {
		t.Errorf("expected no error when writing scratch, but got %s", err.Error())
		return
	}

	if !strings.Contains(path, "test.user@example.com/") {
		t.Error("Unexpected secret path")
		return
	}

	exists, err := s.secretExists(secretMountScratch + "/" + path)
	if err != nil {
		t.Errorf("expected no error when checking secret's existance, but got %s", err.Error())
		return
	}

	if exists == false {
		t.Error("expected to find the secret")
		return
	}
}

func TestVaultClientDeleteSecret(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "http://0.0.0.0:8200", "vault address")
	contextSet.String("vault-token", "1234567890", "vault token")
	contextSet.String("jwt-token", "", "jwt token")
	contextSet.String("vault-role", "youshallnotpass-demo", "vault role")
	contextSet.String("vault-root", "cicd", "vault root")
	contextSet.String("ci-project-path", "youshallnotpass/demo", "ci project path")
	context := cli.NewContext(nil, contextSet, nil)

	secretMountScratch := "cicd/youshallnotpass/demo/scratch"

	s, err := InitVaultClient(context)
	if err != nil {
		t.Errorf("Expected no error when initializing vault client, but got %s", err.Error())
		return
	}

	path, err := s.WriteScratch(0, "test.user@example.com")
	if err != nil {
		t.Errorf("expected no error when writing scratch, but got %s", err.Error())
		return
	}

	if !strings.Contains(path, "test.user@example.com/") {
		t.Error("Unexpected secret path")
		return
	}

	err = s.deleteSecret(secretMountScratch + "/" + path)
	if err != nil {
		t.Errorf("expected no error when deleting scratch, but got %s", err.Error())
		return
	}
}

// Instead of running the actual function I'm going to use hashicorp-related calls in the userMFA function
func TestVaultClientPerformMFA(t *testing.T) {
	contextSet := flag.NewFlagSet("test", 0)
	contextSet.String("vault-addr", "http://0.0.0.0:8200", "vault address")
	contextSet.String("vault-token", "1234567890", "vault token")
	contextSet.String("jwt-token", "", "jwt token")
	contextSet.String("vault-role", "youshallnotpass-demo", "vault role")
	contextSet.String("vault-root", "cicd", "vault root")
	contextSet.String("ci-project-path", "youshallnotpass/demo", "ci project path")
	context := cli.NewContext(nil, contextSet, nil)

	vaultRoot := "cicd"
	ciProjectPath := "youshallnotpass/demo"
	ciPipelineId := 0

	secretMountScratch := vaultRoot + "/" + ciProjectPath + "/" + "scratch"

	ciUserEmail := "test.user@example.com"

	s, err := InitVaultClient(context)
	if err != nil {
		t.Errorf("Expected no error when initializing vault client, but got %s", err.Error())
		return
	}

	secretPath, err := s.WriteScratch(ciPipelineId, ciUserEmail)
	if err != nil {
		t.Errorf("Expected to write scratch correctly: %s", err.Error())
		return
	}

	status, err := s.secretExists(secretMountScratch + "/" + secretPath)
	if !status || err != nil {
		t.Errorf("Expected secret to exist: %s", err.Error())
		return
	}

	err = s.deleteSecret(secretMountScratch + "/" + secretPath)
	if err != nil {
		t.Errorf("Expected to delete secret: %s", err.Error())
		return
	}
}
