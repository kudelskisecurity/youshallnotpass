#!/usr/bin/env bash
# shellcheck disable=SC1091
# shellcheck disable=SC2143

runTest="all"
if [[ -n $1 ]]; then
    runTest=$(echo "$1" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-zA-Z]//g')
fi

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" > /dev/null 2>&1 && pwd)"

source "${currentDir}/../colors.sh"

# build the docker images
if [[ "${runTest}" == "build" ]]; then
    echo -e "${GREEN}Building Integration Tests${NC}"
    docker compose -f "${currentDir}/vaultclient/hashicorp/docker-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_mfa/runner-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_mfa/vault-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_auth_timeout/runner-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_auth_timeout/vault-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_automatic/runner-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_automatic/vault-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_bash/runner-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_bash/vault-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_fail/runner-compose.yml" build
    docker compose -f "${currentDir}/hashicorp_gitlab_fail/vault-compose.yml" build
    echo -e "${BLUE}Finished Building Integration Tests${NC}"
fi

# run hashicorp client tests
if [[ "${runTest}" == "all" || "${runTest}" == "hashicorpclient" || "${runTest}" == "hashicorp" ]]; then
    echo -e "${GREEN}Testing Hashicorp Vault Client${NC}"
    docker compose -f "${currentDir}/vaultclient/hashicorp/docker-compose.yml" up -d >> /dev/null 2>&1

    # wait until vault initialization is complete before we run the hasicorp tests
    sleep 0.2
    while true; do
        if [[ -z $(docker compose -f "${currentDir}/vaultclient/hashicorp/docker-compose.yml" ps | grep "vault_init") ]]; then
            break
        fi
        sleep 0.2
    done

    # run hashicorp tests
    export HASHICORP_TEST_PASSED
    export VAULT_ADDR='http://0.0.0.0:8200'
    if go test "${currentDir}/../../pkg/vaultclient/hashicorpclient"; then
        echo -e "${BLUE}Hashicorp Client Tests PASSED${NC}"
        HASHICORP_TEST_PASSED=true
    else
        echo -e "${RED}Hashicorp Client Tests FAILED${NC}"
        HASHICORP_TEST_PASSED=false
    fi
    docker compose -f "${currentDir}/vaultclient/hashicorp/docker-compose.yml" down >> /dev/null 2>&1
    docker kill hashicorp-vault-1 > /dev/null 2>&1
    echo

    if [[ $HASHICORP_TEST_PASSED == false ]]; then
        exit 1
    fi
fi

# run mattermost client tests
if [[ "${runTest}" == "all" || "${runTest}" == "mattermostclient" || "${runTest}" == "mattermost" ]]; then
    echo -e "${GREEN}Testing Mattermost Client${NC}"
    docker run --name mattermost-preview -d --publish 8065:8065 mattermost/mattermost-preview > /dev/null 2>&1

    echo -e "${GREEN}Letting the Mattermost Docker Image Warm Up / Initialize ${NC}"

    sleep 5
    while true; do
        if curl -i -s http://localhost:8065 > /dev/null 2>&1; then
            echo "Mattermost Docker Host Has Awoken"
            break
        fi
        sleep 5
    done
    
    # run mattermost tests
    export MATTERMOST_TEST_PASSED
    if go test "${currentDir}/../../pkg/loggerclient/mattermostclient"; then
        echo -e "${BLUE}Mattermost Client Tests PASSED${NC}"
        MATTERMOST_TEST_PASSED=true
    else
        echo -e "${RED}Mattermost Client Tests FAILED${NC}"
        MATTERMOST_TEST_PASSED=false
    fi

    docker kill mattermost-preview > /dev/null 2>&1
    docker rm mattermost-preview > /dev/null 2>&1
    echo

    if [[ $MATTERMOST_TEST_PASSED == false ]]; then
        exit 1
    fi
fi

# run hashicorp gitlab automatic integration test
if [[ "${runTest}" == "all" || \
      "${runTest}" == "hashicorp" || \
      "${runTest}" == "gitlab" || \
      "${runTest}" == "automatic" || \
      "${runTest}" == "hashicorp gitlab automatic" || \
      "${runTest}" == "hashicorpgitlabautomatic" ]]; then
    echo -e "${GREEN}Testing Hashicorp-GitLab Automatic Workflow${NC}"
    
    # run vault client
    docker compose -f "${currentDir}/hashicorp_gitlab_automatic/vault-compose.yml" up -d

    # wait for vault to be initialized
    sleep 0.2
    while true; do
        if [[ -z $(docker compose -f "${currentDir}/hashicorp_gitlab_automatic/vault-compose.yml" ps | grep "vault_init") ]]; then
            break
        fi
        sleep 0.2
    done

    # run gitlab runner
    export HASHICORP_GITLAB_AUTOMATIC_PASSED
    if docker compose -f "${currentDir}/hashicorp_gitlab_automatic/runner-compose.yml" up --exit-code-from gitlab_runner; then
        HASHICORP_GITLAB_AUTOMATIC_PASSED=true
    else
        HASHICORP_GITLAB_AUTOMATIC_PASSED=false
    fi

    # shutdown the docker containers
    docker compose -f "${currentDir}/hashicorp_gitlab_automatic/vault-compose.yml" down
    docker compose -f "${currentDir}/hashicorp_gitlab_automatic/runner-compose.yml" down

    if [[ $HASHICORP_GITLAB_AUTOMATIC_PASSED == true ]]; then
        echo -e "\n\n${BLUE}Hashicorp Gitlab Automatic CI/CD Task Test SUCCEEDED${NC}"
    else
        echo -e "\n\n${RED}Hashicorp Gitlab Automatic CI/CD Task Test FAILED${NC}"
        exit 1
    fi
fi

# run hashicorp gitlab failure integration test
if [[ "${runTest}" == "all" || \
      "${runTest}" == "hashicorp" || \
      "${runTest}" == "gitlab" || \
      "${runTest}" == "fail" || \
      "${runTest}" == "hashicorp gitlab fail" || \
      "${runTest}" == "hashicorpgitlabfail" ]]; then
    echo -e "${GREEN}Testing Hashicorp-Gitlab Failure Workflow${NC}"

    # run vault client
    docker compose -f "${currentDir}/hashicorp_gitlab_fail/vault-compose.yml" up -d

    # wait for vault to be initialized
    sleep 0.2
    while true; do
        if [[ -z $(docker compose -f "${currentDir}/hashicorp_gitlab_fail/vault-compose.yml" ps | grep "vault_init") ]]; then
            break
        fi
        sleep 0.2
    done

    # run gitlab runner
    export HASHICORP_GITLAB_FAIL_PASSED
    if docker compose -f "${currentDir}/hashicorp_gitlab_fail/runner-compose.yml" up --exit-code-from gitlab_runner; then
        HASHICORP_GITLAB_FAIL_PASSED=false
    else
        HASHICORP_GITLAB_FAIL_PASSED=true
    fi

    # shutdown the docker containers
    docker compose -f "${currentDir}/hashicorp_gitlab_fail/vault-compose.yml" down
    docker compose -f "${currentDir}/hashicorp_gitlab_fail/runner-compose.yml" down

    if [[ $HASHICORP_GITLAB_FAIL_PASSED == true ]]; then
        echo -e "\n\n${BLUE}Hashicorp Gitlab Failure Workflow Test SUCCEEDED${NC}"
    else
        echo -e "\n\n${RED}Hashicorp Gitlab Failure Workflow Test FAILED${NC}"
        exit 1
    fi
fi

# run hashicorp gitlab auth integration test
if [[ "${runTest}" == "all" || \
      "${runTest}" == "hashicorp" || \
      "${runTest}" == "gitlab" || \
      "${runTest}" == "auth" || \
      "${runTest}" == "hashicorp gitlab auth" || \
      "${runTest}" == "hashicorpgitlabauth" ]]; then
    echo -e "${GREEN}Testing Hashicorp-Gitlab Auth Workflow${NC}"

    # run vault client
    docker compose -f "${currentDir}/hashicorp_gitlab_mfa/vault-compose.yml" up -d

    # wait for vault to be initialized
    sleep 0.2
    while true; do
        if [[ -z $(docker compose -f "${currentDir}/hashicorp_gitlab_mfa/vault-compose.yml" ps | grep "vault_init") ]]; then
            break
        fi
        sleep 0.2
    done

    # run gitlab runner
    export HASHICORP_GITLAB_AUTH_PASSED
    if docker compose -f "${currentDir}/hashicorp_gitlab_mfa/runner-compose.yml" up --exit-code-from gitlab_runner; then
        HASHICORP_GITLAB_AUTH_PASSED=true
    else
        HASHICORP_GITLAB_AUTH_PASSED=false
    fi

    # shutdown the docker containers
    docker compose -f "${currentDir}/hashicorp_gitlab_mfa/vault-compose.yml" down
    docker compose -f "${currentDir}/hashicorp_gitlab_mfa/runner-compose.yml" down


    if [[ $HASHICORP_GITLAB_AUTH_PASSED == true ]]; then
        echo -e "\n\n${BLUE}Hashicorp Gitlab Auth Workflow Test SUCCEEDED${NC}"
    else
        echo -e "\n\n${RED}Hashicorp Gitlab Auth Workflow Test FAILED${NC}"
        exit 1
    fi
fi

# run hashicorp gitlab auth timeout integration test
if [[ "${runTest}" == "all" || \
      "${runTest}" == "hashicorp" || \
      "${runTest}" == "gitlab" || \
      "${runTest}" == "timeout" || \
      "${runTest}" == "hashicorp gitlab timeout" || \
      "${runTest}" == "hashicorpgitlabtimeout" ]]; then
    echo -e "${GREEN}Testing Hashicorp-Gitlab Auth Timeout Workflow${NC}"

    # run vault client
    docker compose -f "${currentDir}/hashicorp_gitlab_auth_timeout/vault-compose.yml" up -d

    # wait for vault to be initialized
    sleep 0.2
    while true; do
        if [[ -z $(docker compose -f "${currentDir}/hashicorp_gitlab_auth_timeout/vault-compose.yml" ps | grep "vault_init") ]]; then
            break
        fi
        sleep 0.2
    done

    # run gitlab runner
    export HASHICORP_GITLAB_TIMEOUT_PASSED
    if docker compose -f "${currentDir}/hashicorp_gitlab_auth_timeout/runner-compose.yml" up --exit-code-from gitlab_runner; then
        HASHICORP_GITLAB_TIMEOUT_PASSED=false
    else
        HASHICORP_GITLAB_TIMEOUT_PASSED=true
    fi

    # shutdown the docker containers
    docker compose -f "${currentDir}/hashicorp_gitlab_auth_timeout/vault-compose.yml" down
    docker compose -f "${currentDir}/hashicorp_gitlab_auth_timeout/runner-compose.yml" down

    if [[ $HASHICORP_GITLAB_TIMEOUT_PASSED == true ]]; then
        echo -e "\n\n${BLUE}Hashicorp Gitlab Auth Timeout Workflow Test SUCCEEDED${NC}"
    else
        echo -e "\n\n${RED}Hashicorp Gitlab Auth Timeout Workflow Test FAILED${NC}"
        exit 1
    fi
fi

if [[ "${runTest}" == "all" || \
      "${runTest}" == "hashicorp" || \
      "${runTest}" == "gitlab" || \
      "${runTest}" == "bash" || \
      "${runTest}" == "hashicorp gitlab bash" || \
      "${runTest}" == "hashicorpgitlabbash" ]]; then
    echo -e "${GREEN}Testing Hashicorp-GitLab Bash Script Workflow${NC}"

    # run vault client
    docker compose -f "${currentDir}/hashicorp_gitlab_bash/vault-compose.yml" up -d

    # wait for vault to be initialized
    sleep 0.2
    while true; do
        if [[ -z $(docker compose -f "${currentDir}/hashicorp_gitlab_bash/vault-compose.yml" ps | grep "vault_init") ]]; then
            break
        fi
        sleep 0.2
    done

    # run gitlab runner
    export HASHICORP_GITLAB_BASH_PASSED
    if docker compose -f "${currentDir}/hashicorp_gitlab_bash/runner-compose.yml" up --exit-code-from gitlab_runner; then
        HASHICORP_GITLAB_BASH_PASSED=true
    else
        HASHICORP_GITLAB_BASH_PASSED=false
    fi

    # shutdown the docker containers
    docker compose -f "${currentDir}/hashicorp_gitlab_bash/vault-compose.yml" down
    docker compose -f "${currentDir}/hashicorp_gitlab_bash/runner-compose.yml" down

    if [[ $HASHICORP_GITLAB_BASH_PASSED == true ]]; then
        echo -e "\n\n${BLUE}Hashicorp Gitlab Bash Script CI/CD Task Test SUCCEEDED${NC}"
    else
        echo -e "\n\n${RED}Hashicorp Gitlab Bash Script CI/CD Task Test FAILED${NC}"
        exit 1
    fi
fi

exit 0