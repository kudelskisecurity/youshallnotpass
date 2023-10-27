package scriptcleanerparser

import (
	"reflect"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/scriptcleanerclient/githubcleanup"
	"github.com/kudelskisecurity/youshallnotpass/pkg/scriptcleanerclient/gitlabcleanup"
)

func TestScriptCleanerParser(t *testing.T) {
	scriptCleanerParserTests := []struct {
		name            string
		platform        string
		expectedCleaner ScriptCleaner
		error           error
	}{
		{
			name:            "parse gitlab cleaner",
			platform:        "gitlab",
			expectedCleaner: &gitlabcleanup.GitLabCleaner{},
			error:           nil,
		},
		{
			name:            "parse github cleaner",
			platform:        "github",
			expectedCleaner: &githubcleanup.GitHubCleaner{},
			error:           nil,
		},
		{
			name:            "fail to parse unknown cleaner client",
			platform:        "dinosaurs",
			expectedCleaner: nil,
			error:           ErrUnknownCICDPlatform,
		},
	}

	for testNum, test := range scriptCleanerParserTests {
		client, err := ParseCleaner(test.platform)
		if err != test.error {
			t.Errorf("\n%d) errors differ -\nexpected: (%+v)\ngot: (%+v)", testNum, test.error, err)
		}

		if reflect.TypeOf(client) != reflect.TypeOf(test.expectedCleaner) {
			t.Errorf("\n%d different client types found -\nexpected: (%s)\ngot: (%s)", testNum, reflect.TypeOf(test.expectedCleaner), reflect.TypeOf(client))
		}
	}
}
