#!/usr/bin/env bash
# shellcheck disable=SC1091
# shellcheck disable=SC2143

runTask="all"
if [[ -n $1 ]]; then
    runTask=$(echo "$1" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-zA-Z]//g')
fi

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" > /dev/null 2>&1 && pwd)"

# get the various color codes from colors.sh
source "${currentDir}/../colors.sh"

# run checkparser tests
if [[ "${runTask}" == "all" || "${runTask}" == "checkparser" ]]; then
    echo -e "${GREEN}Testing CheckParser Client${NC}"
    if go test "${currentDir}/../../pkg/checkparser"; then
        echo -e "${BLUE}CheckParser Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}CheckParser Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run checks tests
if [[ "${runTask}" == "all" || "${runTask}" == "checks" ]]; then
    echo -e "${GREEN}Testing Checks Client${NC}"
    if go test "${currentDir}/../../pkg/checks"; then
        echo -e "${BLUE}Checks Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}Checks Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run imageHashCheck Tests
if [[ "${runTask}" == "all" || "${runTask}" == "imageHash" ]]; then
    echo -e "${GREEN}Testing Image Hash Check Client${NC}"
    if go test "${currentDir}/../../pkg/checks/imagehash"; then
        echo -e "${BLUE}Image Hash Check Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}Image Hash Check Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run scriptHashCheck Tests
if [[ "${runTask}" == "all" || "${runTask}" == "scriptHash" ]]; then
    echo -e "${GREEN}Testing Script Hash Check Client${NC}"
    if go test "${currentDir}/../../pkg/checks/scripthash"; then
        echo -e "${BLUE}Script Hash Check Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}Script Hash Check Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run mfaRequiredCheck Tests
if [[ "${runTask}" == "all" || "${runTask}" == "mfaRequired" ]]; then
    echo -e "${GREEN}Testing Mfa Required Check Client${NC}"
    if go test "${currentDir}/../../pkg/checks/scripthash"; then
        echo -e "${BLUE}Mfa Required Check Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}Mfa Required Check Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run dateTimeCheck Tests
if [[ "${runTask}" == "all" || "${runTask}" == "datetime" ]]; then
    echo -e "${GREEN}Testing Date Time Check Client${NC}"
    if go test "${currentDir}/../../pkg/checks/datetime"; then
        echo -e "${BLUE}Date Time Check Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}Date Time Check Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run config Tests
if [[ "${runTask}" == "all" || "${runTask}" == "config" ]]; then
    echo -e "${GREEN}Testing Config Client${NC}"
    if go test "${currentDir}/../../pkg/config"; then
        echo -e "${BLUE}Config Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}Config Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run scriptcleanupparser tests
if [[ "${runTask}" == "all" || "${runTask}" == "scriptcleanupparser" ]]; then
    echo -e "${GREEN}Testing Script Cleanup Parser${NC}"
    if go test "${currentDir}/../../pkg/scriptcleanerclient"; then
        echo -e "${BLUE}Script Cleanup Parser Tests PASSED${NC}\n"
    else
        echo -e "${RED}Script Cleanup Parser Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run gitlab cleanup tests
if [[ "${runTask}" == "all" || "${runTask}" == "gitlabcleanup" ]]; then
    echo -e "${GREEN}Testing GitLab Cleanup Client${NC}"
    if go test "${currentDir}/../../pkg/scriptcleanerclient/gitlabcleanup"; then
        echo -e "${BLUE}GitLab Cleanup Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}GitLab Cleanup Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run github cleanup tests
if [[ "${runTask}" == "all" || "${runTask}" == "githubcleanup" ]]; then
    echo -e "${GREEN}Testing GitHub Cleanup Client${NC}"
    if go test "${currentDir}/../../pkg/scriptcleanerclient/githubcleanup"; then
        echo -e "${BLUE}GitHub Cleanup Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}GitHub Cleanup Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run vault client tests
if [[ "${runTask}" == "all" || "${runTask}" == "gitlabcleanup" ]]; then
    echo -e "${GREEN}Testing Vault Client${NC}"
    if go test "${currentDir}/../../pkg/vaultclient"; then
        echo -e "${BLUE}Vault Client Tests PASSED${NC}\n"
    else
        echo -e "${RED}Vault Client Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run whitelist tests
if [[ "${runTask}" == "all" || "${runTask}" == "whitelist" ]]; then
    echo -e "${GREEN}Testing Whitelist${NC}"
    if go test "${currentDir}/../../pkg/whitelist"; then
        echo -e "${BLUE}Whitelist Tests PASSED${NC}\n"
    else
        echo -e "${RED}Whitelist Tests FAILED${NC}\n"
        exit 1
    fi
fi

exit 0