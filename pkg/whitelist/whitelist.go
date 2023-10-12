package whitelist

import (
	"errors"
	"strings"
)

var (
	ErrShaNotPresent       = errors.New("@sha256 not provided in image name")
	ErrShaNotValidSha      = errors.New("@Sha 256 not valid sha256. Expected: name:tag@sha256:<sha>")
	ErrInvalidWhitelistJob = errors.New("invalid script name for whitelist script. Expected: <jobName>@sha256:<sha>")
)

type Whitelist struct {
	AllowedImages  []string `json:"allowed_images"`
	AllowedScripts []string `json:"allowed_scripts"`
}

func (s Whitelist) ContainsImage(image string) (bool, error) {
	imageSha, err := getImageSha(image)
	if err != nil {
		return false, err
	}

	for _, v := range s.AllowedImages {
		sha256, _ := getImageSha(v)
		if sha256 == imageSha {
			return true, nil
		}
	}

	return false, nil
}

// Checks through the ScriptWhitelist to see if our script hs allowed.
func (s Whitelist) ContainsScript(scriptSha string) bool {
	for _, script := range s.AllowedScripts {
		// get script plaintext and hash
		sha256, _ := getScriptSha(script)
		if sha256 == scriptSha {
			return true
		}
	}

	return false
}

func (s Whitelist) ContainsJobName(JobName string) (bool, string) {
	for _, script := range s.AllowedScripts {
		whitelist_JobName, _ := getJobName(script)
		if whitelist_JobName == JobName {
			sha256, err := getScriptSha(script)
			if err != nil {
				continue
			}
			return true, sha256
		}
	}

	return false, ""
}

func (w *Whitelist) AddWhitelist(other Whitelist) {
	w.AllowedImages = append(w.AllowedImages, other.AllowedImages...)
	w.AllowedScripts = append(w.AllowedScripts, other.AllowedScripts...)
}

func getJobName(whitelistScript string) (string, error) {
	scriptParts := strings.Split(whitelistScript, "@")

	if len(scriptParts[0]) == 0 {
		return "", ErrInvalidWhitelistJob
	}

	return scriptParts[0], nil
}

func getScriptSha(whitelistScript string) (string, error) {
	script := strings.Split(whitelistScript, "@")

	// 256 bits / 4 bits per char = 64 characters
	shaLen := 256 / 8

	if len(script) == 1 {
		return "", ErrShaNotPresent
	}

	if len(script[1]) < (len("sha256") + shaLen) {
		return "", ErrShaNotValidSha
	}

	return script[1], nil
}

func getImageSha(dockerImage string) (string, error) {
	image := strings.Split(dockerImage, "@")

	// 256 bites / 4 bits per char = 64 characters
	shaLen := 256 / 4

	if len(image) == 1 {
		return "", ErrShaNotPresent
	}

	if len(image[1]) < (len("sha256") + shaLen) {
		return "", ErrShaNotValidSha
	}

	return image[1], nil
}
