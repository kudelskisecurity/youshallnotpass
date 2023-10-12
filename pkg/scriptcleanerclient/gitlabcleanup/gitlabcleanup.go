package gitlabcleanup

import (
	"regexp"
	"strings"
)

type GitLabCleaner struct{}

func (cleaner *GitLabCleaner) CleanupScript(script string) []string {
	re := regexp.MustCompile(`x1b.32;1m.{1,}?x1b.0;m`)
	matches := re.FindAllStringSubmatch(script, -1)
	var lines []string
	for _, match := range matches {
		for _, subMatch := range match {
			line := strings.TrimPrefix(subMatch, `x1b[32;1m`)
			line = strings.TrimSuffix(line, `x1b[0;m`)
			lines = append(lines, line)
		}
	}

	return lines
}
