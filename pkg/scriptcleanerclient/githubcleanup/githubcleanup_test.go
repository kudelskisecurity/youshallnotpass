package githubcleanup_test

import (
	"reflect"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/scriptcleanerclient/githubcleanup"
)

func TestCleanup(t *testing.T) {
	cleanupScriptTests := []struct {
		name           string
		cleaner        githubcleanup.GitHubCleaner
		script         string
		expectedOutput []string
	}{
		{
			name:    "cleanup GitHub scripts no names",
			cleaner: githubcleanup.GitHubCleaner{},
			script: `
- run: actions/checkout@v3
- run: echo "testing"
- run: echo ${{ vars.VAULT_ADDR }}`,
			expectedOutput: []string{
				`run: actions/checkout@v3`,
				`run: echo "testing"`,
				`run: echo ${{ vars.VAULT_ADDR }}`,
			},
		},
		{
			name:    "cleanup GitHub scripts with names",
			cleaner: githubcleanup.GitHubCleaner{},
			script: `
- name: Check out repository
  uses: actions/checkout@v3
- run: echo "testing"
- run: echo ${{ vars.VAULT_ADDR }}`,
			expectedOutput: []string{
				`name: Check out repository
  uses: actions/checkout@v3`,
				`run: echo "testing"`,
				`run: echo ${{ vars.VAULT_ADDR }}`,
			},
		},
		{
			name:    "cleanup GitHub scripts with names and `with`",
			cleaner: githubcleanup.GitHubCleaner{},
			script: `
- name: Check out repository
  uses: actions/checkout@v3
- name: use Node.js
  uses: actions/setup-node@v1
  with:
    node-version: '18.x'`,
			expectedOutput: []string{
				`name: Check out repository
  uses: actions/checkout@v3`,
				`name: use Node.js
  uses: actions/setup-node@v1
  with:
    node-version: '18.x'`,
			},
		},
	}

	for testNum, test := range cleanupScriptTests {
		cleanedScript := test.cleaner.CleanupScript(test.script)
		if !reflect.DeepEqual(test.expectedOutput, cleanedScript) {
			t.Errorf("\n%d) cleaned script was not as expected -\nexpected: (%+v)\ngot: (%+v)", testNum, test.expectedOutput, cleanedScript)
		}
	}
}
