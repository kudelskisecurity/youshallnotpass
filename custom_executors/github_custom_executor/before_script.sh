#!/usr/bin/env bash
# shellcheck disable=SC1091

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

# Get Runner Configuration Variables
source "${currentDir}/profile.sh"

function generate_certs() {
    echo "Generating Pub/Private keys for JWT tokens"
    openssl genrsa -out "${currentDir}/certs/private-key.pem" 3072
    openssl rsa -in "${currentDir}/certs/private-key.pem" -pubout -out "${currentDir}/certs/public-key.pem"
    echo "0" > "${currentDir}/certs/iteration.txt"
}

# check that openssl, jq, and yq are installed
if ! type openssl > /dev/null 2>&1 || ! type jq > /dev/null 2>&1 || ! type yq > /dev/null 2>&1; then
    package_manager=""
    if type brew > /dev/null 2>&1; then
        package_manager="brew install"
    elif type apk > /dev/null 2>&1; then
        package_manager="apk add"
    elif type apt-get > /dev/null 2>&1; then
        package_manager="apt-get install -y"
    elif type dnf > /dev/null 2>&1; then
        package_manager="dnf install"
    elif type yum > /dev/null 2>&1; then
        package_manager="yum install"
    else
        echo "unknown package manager"
    fi

    $package_manager openssl jq yq
fi

# If This is a New Runner We Should Generate a Public-Private Key for Vault Auth
if [[ ! -d "${currentDir}/certs" ]]; then
    mkdir certs
    generate_certs
    echo "Public/Private Key generated in certs folder, please authenticate with vault before continuing"
    exit 1
elif [[ ! -f "${currentDir}/certs/public-key.pem" || ! -f "${currentDir}/certs/private-key.pem" ]]; then
    rm -rf certs
    mkdir certs
    generate_certs
    echo "Public/Private Key generated in certs folder, please authenticate with vault before continuing"
    exit 1
fi

# Check the last modified date of the certs
last_modified=$(date -r "${currentDir}/certs/public-key.pem" +%s)
now=$(date +%s)
declare -i last_modified
declare -i now
last_modified+=15780000 # six months in seconds
if [ $last_modified -le $now ]; then
    echo "it's about time to rotate your public/private key"
fi

# Get Necessary Variables from the environment
export CI_PROJECT_PATH="$GITHUB_REPOSITORY"
export CI_PROJECT_NAMESPACE
CI_PROJECT_NAMESPACE=$(echo "$GITHUB_REPOSITORY" | grep -oE "[a-zA-Z0-9_-]*./" | tr -d "/")
export CI_PIPELINE_ID="$GITHUB_RUN_ID"
export CI_JOB_NAME="$GITHUB_JOB"
export CI_USER_EMAIL="$GITHUB_ACTOR"

# Clone the workflow's repo
# For some reason, the repo is not yet cloned at this stage at the very first run and GITHUB_TOKEN is not available
# Next runs (might?) have the repo locally due to some caching(?) but not the latest version
# sometimes, the repo directory does not contain the .git directory anymore
# sometimes, the git remote -v are erased which breaks for private repos
# TODO: improve me - or improve the github runner itself? https://github.com/actions/runner/issues
# shellcheck disable=SC2086  # not sure how to keep the ls -A simple
if [[ ! -d "${GITHUB_WORKSPACE}" || -z "$(ls -A ${GITHUB_WORKSPACE})" ]]; then
    # set those variables in profile.sh to git clone a private repo
    if [ -n "${GITHUB_USER}" ] && [ -n "${GITHUB_TOKEN}" ]; then
        git clone "https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}" "${GITHUB_WORKSPACE}"
    else
        git clone "${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}" "${GITHUB_WORKSPACE}"
    fi
else
    # directory already exists, force update it
    cd "${GITHUB_WORKSPACE}" || exit 1
    # sometimes, the .git directory does no longer exist...
    if [ ! -d ".git" ]; then
        cd .. || exit 1
        sudo rm -rf "${GITHUB_WORKSPACE}"
        if [ -n "${GITHUB_USER}" ] && [ -n "${GITHUB_TOKEN}" ]; then
            git clone "https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}" "${GITHUB_WORKSPACE}"
        else
            git clone "${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}" "${GITHUB_WORKSPACE}"
        fi
    else
        if [ -n "${GITHUB_USER}" ] && [ -n "${GITHUB_TOKEN}" ]; then
            git remote set-url origin "https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}"
        fi
        git fetch --all
        # fails sometimes...
        # git reset --hard "${GITHUB_REF}"
    fi
fi

# Checkout the current sha
cd "${GITHUB_WORKSPACE}" || exit 1
git config --local advice.detachedHead false
git checkout "${GITHUB_SHA}"
cd "${currentDir}" || exit 1

# Get the Script to be Run
WORKFLOW_FILE_NAME=$(echo "$GITHUB_WORKFLOW_REF" | grep -oE ".github/workflows/[a-zA-Z0-9_-]*\.yml")
WORKFLOW_FILE_NAME=${WORKFLOW_FILE_NAME#".github/workflows/"}
export RUNNER_SCRIPT
RUNNER_SCRIPT=$(yq eval ".jobs.$GITHUB_JOB.steps" "${GITHUB_WORKSPACE}/.github/workflows/$WORKFLOW_FILE_NAME")

# Get the Docker Image to Use (only applies on linux systems)
export CI_JOB_IMAGE=""
CI_JOB_IMAGE=$(yq eval ".jobs.$GITHUB_JOB.container.image" "${GITHUB_WORKSPACE}/.github/workflows/$WORKFLOW_FILE_NAME")
if [ "$CI_JOB_IMAGE" == "null" ]; then
    CI_JOB_IMAGE=$(yq eval ".jobs.$GITHUB_JOB.container" "${GITHUB_WORKSPACE}/.github/workflows/$WORKFLOW_FILE_NAME")
fi

# Generate the jwt token
export ITERATION
ITERATION=$(cat "${currentDir}/certs/iteration.txt")
declare -i ITERATION
ITERATION+=1
echo "$ITERATION" > "${currentDir}/certs/iteration.txt"

export CI_JOB_JWT
CI_JOB_JWT=$("${currentDir}/generate_jwt.sh")

export YOUSHALLNOTPASS_PREVALIDATION_TOKEN
YOUSHALLNOTPASS_PREVALIDATION_TOKEN=$(dd bs=512 if=/dev/urandom count=1 2>/dev/null | LC_ALL=C tr -dc "a-zA-Z0-9" | head -c 50)

# Run youshallnotpass
if [[ "${CI_JOB_IMAGE}" == "null" ]]; then
    "${currentDir}/youshallnotpass" validate --check-type="script" --ci-platform="github" || exit 1
else
    "${currentDir}/youshallnotpass" validate --check-type="all" --ci-platform="github" || exit 1
fi
