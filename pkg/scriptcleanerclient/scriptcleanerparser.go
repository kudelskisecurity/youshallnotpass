package scriptcleanerparser

import (
	"errors"
	"strings"

	"github.com/kudelskisecurity/youshallnotpass/pkg/scriptcleanerclient/githubcleanup"
	"github.com/kudelskisecurity/youshallnotpass/pkg/scriptcleanerclient/gitlabcleanup"
)

var (
	ErrUnknownCICDPlatform = errors.New("unknown CI/CD platform")
)

type ScriptCleaner interface {
	CleanupScript(script string) []string
}

func ParseCleaner(platform string) (ScriptCleaner, error) {
	switch strings.ToLower(platform) {
	case "gitlab":
		return &gitlabcleanup.GitLabCleaner{}, nil
	case "github":
		return &githubcleanup.GitHubCleaner{}, nil
	default:
		return nil, ErrUnknownCICDPlatform
	}
}
