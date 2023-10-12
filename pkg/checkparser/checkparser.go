package checkparser

import (
	"errors"
	"strings"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/checks/datetime"
	"github.com/kudelskisecurity/youshallnotpass/pkg/checks/imagehash"
	"github.com/kudelskisecurity/youshallnotpass/pkg/checks/mfarequired"
	"github.com/kudelskisecurity/youshallnotpass/pkg/checks/scripthash"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
)

var ErrUnknownCheckNameError = errors.New("unknown check name")

func ParseChecks(configs []config.CheckConfig, jobName string, image string, scriptLines []string, checkType string, ciPlatform string) ([]checks.Check, error) {
	var stage uint
	if checkType == "image" {
		stage = checks.ImageCheck
	} else if checkType == "script" {
		stage = checks.ScriptCheck
	} else {
		stage = checks.All
	}

	var performChecks []checks.Check
	for _, config := range configs {
		switch strings.ToLower(config.Name) {
		case "scripthash":
			scriptHashCheck := scripthash.NewScriptHashCheck(config, jobName, scriptLines)
			if scriptHashCheck.IsValidForPlatform(ciPlatform) && scriptHashCheck.IsValidForCheckType(stage) {
				performChecks = append(performChecks, &scriptHashCheck)
			}
		case "imagehash":
			imageHashCheck := imagehash.NewImageHashCheck(config, jobName, image)
			if imageHashCheck.IsValidForPlatform(ciPlatform) && imageHashCheck.IsValidForCheckType(stage) {
				performChecks = append(performChecks, &imageHashCheck)
			}
		case "mfarequired":
			mfaRequiredCheck := mfarequired.NewMfaRequiredCheck(config, jobName)
			if mfaRequiredCheck.IsValidForPlatform(ciPlatform) && mfaRequiredCheck.IsValidForCheckType(stage) {
				performChecks = append(performChecks, &mfaRequiredCheck)
			}
		case "datetimecheck":
			dateTimeCheck := datetime.NewDateTimeCheck(config, jobName)
			if dateTimeCheck.IsValidForPlatform(ciPlatform) && dateTimeCheck.IsValidForCheckType(stage) {
				performChecks = append(performChecks, &dateTimeCheck)
			}
		default:
			return performChecks, ErrUnknownCheckNameError
		}
	}

	return performChecks, nil
}
