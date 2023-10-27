package gitlabcleanup_test

import (
	"reflect"
	"testing"

	"github.com/kudelskisecurity/youshallnotpass/pkg/scriptcleanerclient/gitlabcleanup"
)

func TestGitlabCleanupScript(t *testing.T) {
	cleanupTests := []struct {
		name                  string
		cleaner               gitlabcleanup.GitLabCleaner
		script                string
		expectedCleanedScript []string
	}{
		{
			name:    "cleanup one line gitlab script",
			cleaner: gitlabcleanup.GitLabCleaner{},
			script:  `x1b[32;1m$ echo 'this is the automated job and should be run automatically'x1b[0;m'necho 'this is the automated job and should be run automatically'necho $'x1b[32;1m`,
			expectedCleanedScript: []string{
				`$ echo 'this is the automated job and should be run automatically'`,
			},
		},
		{
			name:    "cleanup multi line gitlab script",
			cleaner: gitlabcleanup.GitLabCleaner{},
			script:  `x1b[32;1m$ echo 'this is the automated job and should be run automatically'x1b[0;m'necho 'this is the automated job and should be run automatically'necho $'x1b[32;1m$ echo "Testing this script"x1b[0;m'necho "Testing this script"necho $'x1b[32;1m$ echo "I just need a few lines to split"x1b[0;m'necho "I just need a few lines to split"necho $'x1b[32;1m$ echo "end"x1b[0;m'necho "end"n'`,
			expectedCleanedScript: []string{
				`$ echo 'this is the automated job and should be run automatically'`,
				`$ echo "Testing this script"`,
				`$ echo "I just need a few lines to split"`,
				`$ echo "end"`,
			},
		},
		{
			name:    "cleanup bash script",
			cleaner: gitlabcleanup.GitLabCleaner{},
			script:  `x1b[32;1m$ /gitrepo/test.shx1b[0;m'n/gitrepo/test.shn'x1b[32;1m/gitrepo/test.shx1b[0;mx1b[32;1m#!/bin/bash;;echo "TEST SCRIPT";;x1b[0;m`,
			expectedCleanedScript: []string{
				`$ /gitrepo/test.sh`,
				`/gitrepo/test.sh`,
				`#!/bin/bash;;echo "TEST SCRIPT";;`,
			},
		},
		{
			name:    "cleanup combined bash + script",
			cleaner: gitlabcleanup.GitLabCleaner{},
			script:  `x1b[32;1m$ echo "this is the script job"x1b[0;m'necho "this is the script job"necho $'x1b[32;1m$ /gitrepo/test.shx1b[0;m'n/gitrepo/test.shnecho $'x1b[32;1m$ echo "this is the end of the script job"x1b[0;m'necho "this is the end of the script job"n'x1b[32;1m/gitrepo/test.shx1b[0;mx1b[32;1m#!/bin/bash;;echo "TEST SCRIPT";;x1b[0;m`,
			expectedCleanedScript: []string{
				`$ echo "this is the script job"`,
				`$ /gitrepo/test.sh`,
				`$ echo "this is the end of the script job"`,
				`/gitrepo/test.sh`,
				`#!/bin/bash;;echo "TEST SCRIPT";;`,
			},
		},
	}

	for testNum, test := range cleanupTests {
		cleanedScript := test.cleaner.CleanupScript(test.script)
		if !reflect.DeepEqual(test.expectedCleanedScript, cleanedScript) {
			t.Errorf("\n%d) cleaned script was not expected -\nexpected: %+v\ngot: %+v", testNum, test.expectedCleanedScript, cleanedScript)
		}
	}
}
