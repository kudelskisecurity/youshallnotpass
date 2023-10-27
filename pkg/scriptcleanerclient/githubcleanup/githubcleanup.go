package githubcleanup

import "strings"

type GitHubCleaner struct{}

func (cleaner *GitHubCleaner) CleanupScript(script string) []string {
	/*
		GitHub steps start of the form which is almost valid yaml, we're
		going to remove the "- "s
		- name: Check out repository
		  uses: actions/checkout@v3
		- run: echo "testing"
		- run: echo ${{ vars.VAULT_ADDR }}

		Which is parsed to the form:
		`name: Check out repository
		 	uses: actions/checkout@v3`,
		`run: echo "testing"`,
		`run: echo ${{ vars.VAULT_ADDR }}`
	*/
	splitScript := strings.Split(script, "- ")
	splitScript = splitScript[1:]
	var scriptLines []string
	for _, section := range splitScript {
		section = strings.TrimSpace(section)
		section = strings.ReplaceAll(section, "\t", "    ")
		scriptLines = append(scriptLines, section)
	}
	return scriptLines
}
