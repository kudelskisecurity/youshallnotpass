package mattermostclient_test

import (
	"context"
	"strings"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/loggerclient/mattermostclient"
	"github.com/mattermost/mattermost/server/public/model"
)

var (
	client *mattermostclient.MattermostClient = nil
)

// This is more integration testing than unit testing because the hashicorp vault
// client is pretty much known to work we're just wrapping it

func initialize() error {
	if client == nil {
		ctx := context.Background()
		mmClient := model.NewAPIv4Client("http://localhost:8065")

		user, _, err := mmClient.CreateUser(ctx, &model.User{
			Username: "youshallnotpass-bot",
			Email:    "test@gmail.com",
			Password: "1234567890",
		})
		if err != nil {
			return err
		}

		_, _, err = mmClient.LoginById(ctx, user.Id, "1234567890")
		if err != nil {
			return err
		}

		team, _, err := mmClient.CreateTeam(ctx, &model.Team{
			Name:        "youshallnotpass",
			DisplayName: "youshallnotpass",
			Type:        "O",
		})
		if err != nil {
			return err
		}

		youshallnotpassChannel, _, err := mmClient.CreateChannel(ctx, &model.Channel{
			TeamId:      team.Id,
			Name:        "youshallnotpass",
			DisplayName: "youshallnotpass",
			Type:        "O",
		})
		if err != nil {
			return err
		}

		mattermostClient, err := mattermostclient.NewMattermostClient("http://localhost:8065", mmClient.AuthToken, youshallnotpassChannel.Id, "test-job", "http://127.0.0.1:8200", "cicd", "test")
		if err != nil {
			return err
		}

		client = mattermostClient
	}
	return nil
}

func TestNewMattermostClient(t *testing.T) {
	err := initialize()

	if err != nil {
		t.Errorf("Expected no err, but got %s", err.Error())
		return
	}
}

func TestLogCheckResults(t *testing.T) {
	err := initialize()

	if err != nil {
		t.Errorf("unexpected err: %s", err.Error())
		return
	}

	checkResults := []checks.CheckResult{
		{
			Name:    "imageHash",
			Version: "1.0.0",
			Error:   nil,
			Abort:   true,
			Mfa:     false,
			Details: "New Image Detected",
		},
		{
			Name:    "scriptHash",
			Version: "1.0.0",
			Error:   nil,
			Abort:   false,
			Mfa:     false,
			Details: "Script Hash Check Succeeded",
		},
		{
			Name:    "mfaRequired",
			Version: "1.0.0",
			Error:   nil,
			Abort:   false,
			Mfa:     true,
			Details: "Mfa required for job",
		},
	}

	err = client.LogCheckResults(checkResults)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestLogPrevalidated(t *testing.T) {
	err := initialize()

	if err != nil {
		t.Errorf("Expected no err, but got %s", err.Error())
		return
	}

	err = client.LogPrevalidated()

	if err != nil {
		t.Errorf("Expected no err, but got %s", err.Error())
		return
	}
}

func TestLogSuccessfulExecution(t *testing.T) {
	err := initialize()

	if err != nil {
		t.Errorf("Expected no err, but got %s", err.Error())
		return
	}

	err = client.LogSuccessfulExecution()

	if err != nil {
		t.Errorf("Expected no err, but got %s", err.Error())
		return
	}
}

func TestLogFailedExecution(t *testing.T) {
	err := initialize()

	if err != nil {
		t.Errorf("Expected no err, but got %s", err.Error())
		return
	}

	err = client.LogFailedExecution("test")

	if !strings.Contains(err.Error(), "test") {
		t.Errorf("Expected error to contain test, instead got %s", err.Error())
		return
	}
}

func TestSendMFAInstructions(t *testing.T) {
	err := initialize()

	if err != nil {
		t.Errorf("Expected no err, but got %s", err.Error())
		return
	}

	err = client.SendMFAInstructions("delete the scratch code")

	if err != nil {
		t.Errorf("Expected no err, but got %s", err.Error())
		return
	}
}
