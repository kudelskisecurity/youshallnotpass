package imagehash

import (
	"strings"
	"sync"

	"github.com/kudelskisecurity/youshallnotpass/pkg/checks"
	"github.com/kudelskisecurity/youshallnotpass/pkg/config"
	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

var (
	ImageHashCheckFailedMfaRequired = "Unknown Image MFA Required"
	ImageHashCheckFailedAbort       = "Unknown Image Aborting"
	ImageHashCheckSuccess           = "Successful Image Hash Check"
	ImageHashCheckError             = "---ERROR---"
)

type ImageHashCheck struct {
	jobName     string
	abortOnFail bool
	mfaOnFail   bool
	image       string
}

func NewImageHashCheck(config config.CheckConfig, jobName string, image string) ImageHashCheck {
	abortOnFail, exists := config.Options["abortOnFail"].(bool)
	if !exists {
		abortOnFail = false
	}

	mfaOnFail, exists := config.Options["mfaOnFail"].(bool)
	if !exists {
		mfaOnFail = false
	}

	return ImageHashCheck{
		abortOnFail: abortOnFail,
		mfaOnFail:   mfaOnFail,
		jobName:     jobName,
		image:       image,
	}
}

func (check *ImageHashCheck) Check(channel chan<- checks.CheckResult, wg *sync.WaitGroup, w whitelist.Whitelist) {
	defer wg.Done()

	result := checks.CheckResult{Name: "Image Hash Check", Version: "1.0.0"}
	details := ""

	foundImg, err := w.ContainsImage(check.image)
	if err != nil {
		result.Error = err
		details = ImageHashCheckError
		if check.abortOnFail {
			result.Abort = true
		} else if check.mfaOnFail {
			result.Mfa = true
		}
		result.Details = details
		channel <- result
		return
	}

	result.Error = nil

	if !foundImg && check.abortOnFail {
		result.Abort = true
		details = ImageHashCheckFailedAbort
	} else if !foundImg && check.mfaOnFail {
		result.Mfa = true
		details = ImageHashCheckFailedMfaRequired
	} else if foundImg {
		details = ImageHashCheckSuccess
	}

	result.Details = strings.Clone(details)

	channel <- result
}

func (check *ImageHashCheck) IsValidForCheckType(checkType uint) bool {
	switch checkType {
	case checks.All:
		return true
	default:
		return checkType == checks.ImageCheck
	}
}

func (check *ImageHashCheck) IsValidForPlatform(ciPlatform string) bool {
	return true
}
