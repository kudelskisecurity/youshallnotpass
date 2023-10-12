package mattermostclient

import (
	"context"
	"fmt"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/mattermost/mattermost/server/public/model"
)

type MattermostClient struct {
	client            *model.Client4
	channelId         string
	postId            string
	ciJobName         string
	vaultExternalAddr string
	vaultRoot         string
	ciProjectPath     string
}

func NewMattermostClient(mattermostUrl string, mattermostToken string, channelId string, ciJobName string, vaultExternalAddr string, vaultRoot string, ciProjectPath string) (*MattermostClient, error) {
	client := model.NewAPIv4Client(mattermostUrl)
	client.SetOAuthToken(mattermostToken)
	ctx := context.Background()

	post, _, err := client.CreatePost(ctx,
		&model.Post{
			ChannelId: channelId,
			Message:   fmt.Sprintf("## %s CI/CD Run", ciJobName),
		},
	)
	if err != nil {
		return nil, err
	}

	return &MattermostClient{
		client:            client,
		channelId:         channelId,
		postId:            post.Id,
		ciJobName:         ciJobName,
		vaultExternalAddr: vaultExternalAddr,
		vaultRoot:         vaultRoot,
		ciProjectPath:     ciProjectPath,
	}, nil
}

func (c *MattermostClient) LogRecoverableError(err error) {
	_, _, _ = c.client.CreatePost(context.Background(),
		&model.Post{
			ChannelId: c.channelId,
			RootId:    c.postId,
			Message:   err.Error(),
		},
	)
}

func (c *MattermostClient) SendMFAInstructions(message string) error {
	ctx := context.Background()
	_, _, err := c.client.CreatePost(ctx,
		&model.Post{
			ChannelId: c.channelId,
			RootId:    c.postId,
			Message:   message,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *MattermostClient) LogCheckResults(results []checks.CheckResult) error {
	message := "---------------------------------------------------------------------\n"
	message += fmt.Sprintf("%s/ui/vault/secrets/%s - %s\n", c.vaultExternalAddr, c.vaultRoot+"/show/"+c.ciProjectPath+"/whitelist", c.ciJobName)
	message += "---------------------------------------------------------------------\n"
	message += "   Name   |   Version   |   Error   |   Abort   |   Mfa   |   Details\n"
	message += "---------------------------------------------------------------------\n"
	for _, result := range results {
		errStr := "nil"
		if result.Error != nil {
			errStr = result.Error.Error()
		}
		message += fmt.Sprintf("%s | %s | %s | %t | %t | %s\n", result.Name, result.Version,
			errStr, result.Abort, result.Mfa, result.Details)
		message += "---------------------------------------------------------------------\n"
	}

	ctx := context.Background()
	_, _, err := c.client.CreatePost(ctx,
		&model.Post{
			ChannelId: c.channelId,
			RootId:    c.postId,
			Message:   message,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *MattermostClient) LogPrevalidated() error {
	message := fmt.Sprintf("✅ CI/CD for %s has been prevalidated", c.ciJobName)
	ctx := context.Background()
	_, _, err := c.client.CreatePost(ctx,
		&model.Post{
			ChannelId: c.channelId,
			RootId:    c.postId,
			Message:   message,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *MattermostClient) LogSuccessfulExecution() error {
	message := fmt.Sprintf("✅ Successful YouShallNotPass check for Job: %s", c.ciJobName)
	ctx := context.Background()
	_, _, err := c.client.CreatePost(ctx,
		&model.Post{
			ChannelId: c.channelId,
			RootId:    c.postId,
			Message:   message,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *MattermostClient) LogFailedExecution(reason string) error {
	message := fmt.Sprintf("❌ Unsuccessful YouShallNotPass check for job: %s\n", c.ciJobName)
	ctx := context.Background()
	_, _, _ = c.client.CreatePost(ctx,
		&model.Post{
			ChannelId: c.channelId,
			RootId:    c.postId,
			Message:   message,
		},
	)

	return fmt.Errorf(reason)
}
