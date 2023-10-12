package hashicorpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/loggerclient"
	"github.com/kudelskisecurity/youshallnotpass/pkg/vaultclient"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
	"github.com/urfave/cli/v2"
)

type HashicorpService struct {
	client             *api.Client
	vaultAddr          string
	vaultExternalAddr  string
	vaultToken         string
	vaultRole          string
	vaultRoot          string
	ciProjectPath      string
	preValidationToken string
}

// Creates a new hashicorp service
func InitVaultClient(context *cli.Context) (*HashicorpService, error) {
	vaultAddr := context.String("vault-addr")
	vaultExternalAddr := context.String("vault-external-addr")
	vaultToken := context.String("vault-token")
	vaultRole := context.String("vault-role")
	vaultRoot := context.String("vault-root")
	ciProjectPath := context.String("ci-project-path")
	preValidationToken := context.String("pre-validation-token")

	if vaultRole == "" {
		vaultRole = strings.ReplaceAll(ciProjectPath, "/", "-")
	}

	jwtToken := context.String("jwt-token")
	loginPath := context.String("vault-login-path")

	if vaultToken == "" && jwtToken == "" {
		err := errors.New("either vaultToken or jwtToken is required")
		return nil, err
	}

	config := &api.Config{
		Address: vaultAddr,
	}

	vaultClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	s := HashicorpService{
		client:             vaultClient,
		vaultAddr:          vaultAddr,
		vaultExternalAddr:  vaultExternalAddr,
		vaultRole:          vaultRole,
		vaultRoot:          vaultRoot,
		ciProjectPath:      ciProjectPath,
		preValidationToken: preValidationToken,
	}

	if jwtToken != "" {
		err = s.getAuthToken(jwtToken, vaultRole, loginPath)
		if err != nil {
			return nil, fmt.Errorf("unable to authenticate to vault: %s", err.Error())
		}
	} else {
		s.vaultToken = vaultToken
	}

	return &s, nil
}

func (s *HashicorpService) GetNamespaceConfig(namespaceConfigMount string) (config.NamespaceConfig, error) {
	vaultRes, err := s.getConfig(namespaceConfigMount)
	if err != nil {
		return config.DefaultNamespaceConfig, fmt.Errorf("unable to access namespace vault config at %s", namespaceConfigMount)
	}

	return config.ParseNamespaceConfig(vaultRes)
}

func (s *HashicorpService) GetProjectConfig(projectConfigMount string) (config.ProjectConfig, error) {
	vaultRes, err := s.getConfig(projectConfigMount)
	if err != nil {
		return config.DefaultProjectConfig, fmt.Errorf("unable to access project vault config at %s", projectConfigMount)
	}

	return config.ParseProjectConfig(vaultRes)
}

func (s *HashicorpService) getConfig(mount string) ([]byte, error) {
	s.client.SetToken(s.vaultToken)

	secret, err := s.client.Logical().Read(mount)
	var vaultRes []byte

	if err != nil {
		return vaultRes, err
	}

	// Run with default config
	if secret == nil {
		return vaultRes, nil
	}

	vaultRes, err = json.Marshal(secret.Data)
	if err != nil {
		return vaultRes, err
	}

	return vaultRes, nil
}

func (s *HashicorpService) ReadWhitelists(namespaceMount string, projectMount string) (whitelist.Whitelist, error) {
	s.client.SetToken(s.vaultToken)

	namespaceWhitelist, err := s.readWhitelist(namespaceMount)
	if err != nil {
		return namespaceWhitelist, fmt.Errorf("unable to read namespace whitelist at %s: %s", namespaceMount, err.Error())
	}

	projectWhitelist, err := s.readWhitelist(projectMount)
	if err != nil {
		return projectWhitelist, fmt.Errorf("unable to read project whitelist at %s: %s", projectMount, err.Error())
	}

	namespaceWhitelist.AddWhitelist(projectWhitelist)

	return namespaceWhitelist, nil
}

// Read the allowed images and scripts from a whitelist mount location
func (s *HashicorpService) readWhitelist(secretPath string) (whitelist.Whitelist, error) {
	s.client.SetToken(s.vaultToken)

	secret, err := s.client.Logical().Read(secretPath)
	whitelist := whitelist.Whitelist{}

	if err != nil {
		return whitelist, err
	}

	if secret == nil {
		return whitelist, fmt.Errorf("invalid secret")
	}

	vaultRes, _ := json.Marshal(secret.Data)
	err = json.Unmarshal(vaultRes, &whitelist)

	if err != nil {
		return whitelist, err
	}

	return whitelist, nil
}

// Write scratch code and make sure it exists after writing
func (s *HashicorpService) WriteScratch(pipelineId int, user string) (string, error) {
	secretMount := s.vaultRoot + "/" + s.ciProjectPath + "/" + "scratch"

	randomString, err := vaultclient.GenerateRandomStringURLSafe(16)
	if err != nil {
		return "", errors.New("an unknown error has occured with string generation")
	}

	s.client.SetToken(s.vaultToken)

	secret := make(map[string]interface{})
	secret["CI/CD pipeline id"] = strconv.Itoa(pipelineId)

	secretPath := user + "/" + randomString
	_, err = s.client.Logical().Write(secretMount+"/"+secretPath, secret)
	if err != nil {
		return "", err
	}

	status, err := s.secretExists(secretMount + "/" + secretPath)
	if !status || err != nil {
		return secretPath, errors.New("  ❌ CI/CD run not authorized, secret not retrievable")
	}

	return secretPath, err
}

// Print to the console the instructions for deleting the scratch code
func (s *HashicorpService) LogMFAInstructions(ciUserEmail string, loggerClient loggerclient.LoggerClient) {
	message := fmt.Sprintf("\nPlease delete the following scratch code to authorize this pipeline run -> %s/ui/vault/secrets/%s\n",
		s.vaultExternalAddr, s.vaultRoot+"/list/"+s.ciProjectPath+"/scratch/"+ciUserEmail)
	_ = loggerClient.SendMFAInstructions(message)
}

// Waits for the user to delete the scratch code in the hashicorp vault
func (s *HashicorpService) WaitForMFA(timeout int, secretPath string) bool {
	secretMount := s.vaultRoot + "/" + s.ciProjectPath + "/" + "scratch"

	vaultCheckIntervalSeconds := 5

	// Every 5 seconds check whether the scratch code has been deleted
	for i := 1; i <= (timeout / vaultCheckIntervalSeconds); i++ {
		status, err := s.secretExists(secretMount + "/" + secretPath)

		if err != nil {
			fmt.Printf("%s", err.Error())
		}

		if !status {
			return true
		}

		time.Sleep(time.Duration(vaultCheckIntervalSeconds) * time.Second)
	}

	return false
}

// If the user failed MFA, delete the scratch code
func (s *HashicorpService) Cleanup(successful bool, secretPath string, checkType string) error {
	secretMount := s.vaultRoot + "/" + s.ciProjectPath + "/" + "scratch"

	if !successful {
		err := s.deleteSecret(secretMount + "/" + secretPath)
		if err != nil {
			return err
		}
	} else if checkType != "image" {
		err := vaultclient.CreateValidationToken(s.preValidationToken)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAuthToken obtains a vaultToken from a JWT token
func (s *HashicorpService) getAuthToken(jwtToken string, role string, loginPath string) error {
	data := map[string]interface{}{
		"jwt":  jwtToken,
		"role": role,
	}

	token, err := s.client.Logical().Write(loginPath, data)
	if err != nil {
		return err
	}
	s.vaultToken = token.Auth.ClientToken
	return nil
}

// Check that a secret (scratch code) exists in the vault
func (s *HashicorpService) secretExists(secretPath string) (bool, error) {
	s.client.SetToken(s.vaultToken)

	secret, err := s.client.Logical().Read(secretPath)

	if err != nil {
		return true, err
	}
	if secret == nil {
		return false, nil
	}

	return true, nil
}

// Delete a secret (scratch code) from the vault
func (s *HashicorpService) deleteSecret(secretPath string) error {
	s.client.SetToken(s.vaultToken)

	_, err := s.client.Logical().Delete(secretPath)

	if err != nil {
		return err
	}

	return nil
}
