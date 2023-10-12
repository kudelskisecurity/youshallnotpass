#!/usr/bin/env bash

testType="all"
if [[ -n $1 ]]; then
    testType=$(echo "$1" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-zA-Z]//g')
fi

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" > /dev/null 2>&1 && pwd)"

# shellcheck disable=SC1091
source "${currentDir}/colors.sh"

# run unit tests
if [[ "${testType}" == "all" || "${testType}" == "unit" ]]; then
    echo -e "${GREEN}Running Unit Tests${NC}"
    # TODO: Print Logo

    if "${currentDir}/unit/test.sh" "$2"; then
        echo -e "${BLUE}Unit Tests PASSED${NC}\n"
    else
        echo -e "${RED}Unit Tests FAILED${NC}\n"
        exit 1
    fi
fi

# run integration tests
if [[ "${testType}" == "all" || "${testType}" == "integration" ]]; then
    echo -e "${GREEN}Running Integration Tests${NC}"
    # TODO: Print Logo

    if "${currentDir}/integration/test.sh" "$2"; then
        echo -e "${BLUE}Integration Tests PASSED${NC}\n"
    else
        echo -e "${RED}Integration Tests FAILED${NC}\n"
        exit 1
    fi
fi

exit 0