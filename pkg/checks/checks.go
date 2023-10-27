package checks

import (
	"sync"

	"github.com/kudelskisecurity/youshallnotpass/pkg/whitelist"
)

const (
	ImageCheck = iota
	ScriptCheck
	All
)

// The Check interface should be implemented for all YouShallNotPass checks.
type Check interface {
	// Check performs whatever logic is necessary to determine whether an execution is
	// allowed before sending the result through a channel and signalling to the wait group
	// that it has completed execution.
	//
	// Example:
	// The ImageHash check checks that the image hash of the current job (obtained on instantiation)
	// exists in the whitelist.  If it does exist, the check creates a CheckResult indicating the check
	// was passed and sends the check through the channel.  Once this result is passed through the channel
	// the wait group is told the check is done.
	Check(chan<- CheckResult, *sync.WaitGroup, whitelist.Whitelist)

	// IsValidForCheckType returns true if a given check is valid for a given check
	// type (i.e. ImageCheck, ScriptCheck, AllCheck)
	//
	// Example:
	// If there is a check named ScriptLintCheck that is only valid for script and all checks, this
	// function will return true if checkType == ImageCheck OR checkType == All.
	IsValidForCheckType(checkType uint) bool

	// IsValidForPlatform returns true if a given check is valid for a given platform
	//
	// Example:
	// If there is a check named ScriptLintCheck that is only valid on GitLab, this function
	// will return true if the ciPlatform == "gitlab"
	IsValidForPlatform(ciPlatform string) bool
}

type CheckResult struct {
	Name    string
	Version string
	Error   error
	Abort   bool
	Mfa     bool
	Details string
}

// Deep Compare Two Check Results (Only Useful for the Testing).
func (result *CheckResult) CompareCheckResult(other CheckResult) bool {
	if result.Name != other.Name {
		return false
	}

	if result.Version != other.Version {
		return false
	}

	if result.Error != other.Error {
		return false
	}

	if result.Abort != other.Abort {
		return false
	}

	if result.Mfa != other.Mfa {
		return false
	}

	if result.Details != other.Details {
		return false
	}

	return true
}
