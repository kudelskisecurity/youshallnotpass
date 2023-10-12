package scripthash

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

var (
	ScriptHashCheckName                 = "Script Hash Check"
	ScriptHashCheckVersion              = "1.0.0"
	ScriptHashCheckEmptyScriptDetails   = "No Script"
	ScriptHashCheckAbortScriptDetails   = "Unknown Script %s@%s Aborting"
	ScriptHashCheckMfaScriptDetails     = "Unknown Script %s@%s MFA Required"
	ScriptHashCheckSuccessDetails       = "Found Script Sha"
	ScriptHashCheckUpdatedScriptDetails = " - CI Job %s has been updated"
)

type ScriptHashCheck struct {
	jobName     string
	abortOnFail bool
	mfaOnFail   bool
	scriptLines []string
}

func NewScriptHashCheck(config config.CheckConfig, jobName string, scriptLines []string) ScriptHashCheck {
	abortOnFail, exists := config.Options["abortOnFail"].(bool)
	if !exists {
		abortOnFail = false
	}

	mfaOnFail, exists := config.Options["mfaOnFail"].(bool)
	if !exists {
		mfaOnFail = false
	}

	return ScriptHashCheck{
		abortOnFail: abortOnFail,
		mfaOnFail:   mfaOnFail,
		jobName:     jobName,
		scriptLines: scriptLines,
	}
}

func (check *ScriptHashCheck) Check(channel chan<- checks.CheckResult, wg *sync.WaitGroup, w whitelist.Whitelist) {
	defer wg.Done()

	result := checks.CheckResult{Name: ScriptHashCheckName, Version: ScriptHashCheckVersion}
	details := ""

	scriptUpdate := false

	scriptSha := check.hashScript()
	if scriptSha == "" {
		details = ScriptHashCheckEmptyScriptDetails
		result.Details = details
		channel <- result
		return
	}

	foundScript := w.ContainsScript(scriptSha)

	if !foundScript {
		scriptUpdate, _ = w.ContainsJobName(check.jobName)
	}

	if !foundScript && check.abortOnFail {
		result.Abort = true
		details = fmt.Sprintf(ScriptHashCheckAbortScriptDetails, check.jobName, scriptSha)
	} else if !foundScript && check.mfaOnFail {
		result.Mfa = true
		details = fmt.Sprintf(ScriptHashCheckMfaScriptDetails, check.jobName, scriptSha)
	} else if foundScript {
		details = ScriptHashCheckSuccessDetails
	}

	if scriptUpdate {
		details += fmt.Sprintf(ScriptHashCheckUpdatedScriptDetails, check.jobName)
	}

	result.Details = strings.Clone(details)

	channel <- result
}

func (check *ScriptHashCheck) hashScript() string {
	script := ""
	for _, line := range check.scriptLines {
		script += line
	}

	if len(script) == 0 {
		return ""
	}

	h := sha256.New()
	h.Write([]byte(script))

	return "sha256:" + base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (check *ScriptHashCheck) IsValidForCheckType(checkType uint) bool {
	switch checkType {
	case checks.All:
		return true
	default:
		return checkType == checks.ScriptCheck
	}
}

func (check *ScriptHashCheck) IsValidForPlatform(ciPlatform string) bool {
	return true
}
