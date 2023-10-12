package mfarequired

import (
	"strings"
	"sync"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

var (
	MfaRequiredCheckName    = "Mfa Required Check"
	MfaRequiredCheckVersion = "1.0.0"
	MfaRequiredCheckDetails = "Mfa Required"
)

type MfaRequiredCheck struct {
	jobName        string
	validCheckType uint
}

func NewMfaRequiredCheck(config config.CheckConfig, jobName string) MfaRequiredCheck {
	validCheckType := checks.All

	checkType, exists := config.Options["checkType"].(string)
	if !exists {
		validCheckType = checks.All
	} else if strings.ToLower(checkType) == "image" {
		validCheckType = checks.ImageCheck
	} else if strings.ToLower(checkType) == "script" {
		validCheckType = checks.ScriptCheck
	}

	return MfaRequiredCheck{
		jobName:        jobName,
		validCheckType: uint(validCheckType),
	}
}

func (check *MfaRequiredCheck) Check(channel chan<- checks.CheckResult, wg *sync.WaitGroup, w whitelist.Whitelist) {
	defer wg.Done()

	result := checks.CheckResult{Name: MfaRequiredCheckName, Version: MfaRequiredCheckVersion, Mfa: true, Details: MfaRequiredCheckDetails}

	channel <- result
}

func (check *MfaRequiredCheck) IsValidForCheckType(checkType uint) bool {
	switch check.validCheckType {
	case checks.All:
		return true
	default:
		switch checkType {
		case checks.All:
			return true
		default:
			return check.validCheckType == checkType
		}
	}
}

func (check *MfaRequiredCheck) IsValidForPlatform(ciPlatform string) bool {
	return true
}
