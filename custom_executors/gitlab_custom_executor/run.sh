#!/usr/bin/env bash
# shellcheck disable=SC1091
# shellcheck disable=SC2129

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

source "${currentDir}/base.sh"

source "${currentDir}/profile.sh"

STEPS_WITH_NO_MFA=(prepare_script get_sources restore_cache download_artifacts archive_cache archive_cache_on_failure upload_artifacts_on_success upload_artifacts_on_failure cleanup_file_variables)

MFA_REQUIRED=true

for i in "${STEPS_WITH_NO_MFA[@]}"
do
    if [ "$i" == "$2" ] ; then
        MFA_REQUIRED=false
    fi
done

GREEN="\x1b[32;1m"
[[ $(grep -Fq "$GREEN" < "$1"; echo $?) == 0 ]] && CONTAINS_SCRIPT=true || CONTAINS_SCRIPT=false

# if user has already MFA'd let them execute a shell
if [ "$MFA_REQUIRED" == true ] && [ "$CONTAINS_SCRIPT" == true ]; then
    # As of writing this, the script began when the green text was, therefore this searches for the first
    # green text occurance and grabs that to the end of file.
    export RUNNER_SCRIPT
    RUNNER_SCRIPT=$(grep -oE "x1b.32;1m.*" < "$1")
    # shellcheck disable=SC1003
    RUNNER_SCRIPT=$(echo "$RUNNER_SCRIPT" | tr -d '\\')

    # Find the scripts
    SCRIPTS=$(echo "$RUNNER_SCRIPT" | grep -oE "(./|/|[a-zA-Z0-9_-])([a-zA-Z0-9_-]*/){0,}.{0,1}[a-zA-Z0-9_-]*\.(sh|py|ps1|rb|js)")

    # Expand the scripts
    for script in $SCRIPTS; do
        script=$(echo "$script" | tr -d '[:space:]')
        if [[ -f $script ]]; then
            ADD_SCRIPT='x1b[32;1m'
            ADD_SCRIPT="${ADD_SCRIPT}${script}"
            ADD_SCRIPT="${ADD_SCRIPT}"'x1b[0;m'
            ADD_SCRIPT="${ADD_SCRIPT}"'x1b[32;1m'$(tr '\n' ';' < "${script}" )'x1b[0;m'
            RUNNER_SCRIPT="${RUNNER_SCRIPT}${ADD_SCRIPT}"
        fi
    done

    # As of writing this, the runner script gives the CI_JOB_NAME environment variable right before
    # CI_JOB_STAGE.
    export CI_JOB_NAME
    CI_JOB_NAME=$(grep -oE "CI_JOB_NAME=.{0,5}[a-zA-Z0-9_-]*" < "$1")
    CI_JOB_NAME=$(echo "$CI_JOB_NAME" | grep -oE "=.*" | grep -oE "[a-zA-Z0-9_-]*")

    # As of writing this, the runner script gives the CI_JOB_STATUS environment variable right before
    # the CI={true/false} variable
    export CI_JOB_STATUS
    CI_JOB_STATUS=$(grep -oE "CI_JOB_STATUS=.{0,5}[a-zA-Z0-9_-]*" < "$1")
    CI_JOB_STATUS=$(echo "$CI_JOB_STATUS" | grep -oE "=.*" | grep -oE "[a-zA-Z0-9_-]*")

    # GitLab has support for before_scripts, scripts and after_scripts.  The before_scripts and scripts
    # Get combined in the file, however the after_scripts seem to be send with the environment variable
    # CI_JOB_STATUS=success, so their related job name will be {CI_JOB_NAME}-after.
    if [ "$CI_JOB_STATUS" == "success" ]; then
        CI_JOB_NAME="$CI_JOB_NAME-after"
    fi

    if [[ -f /usr/local/bin/youshallnotpass ]]; then
        /usr/local/bin/youshallnotpass validate --check-type="script" --ci-platform="gitlab" || exit "$BUILD_FAILURE_EXIT_CODE"
    elif [[ -f ${currentDir}/youshallnotpass ]]; then
        "${currentDir}/youshallnotpass" validate --check-type="script" --ci-platform="gitlab" || exit "$BUILD_FAILURE_EXIT_CODE"
    else
        echo "Could Not Find YouShallNotPass Binary"
        exit "$BUILD_FAILURE_EXIT_CODE"
    fi
fi

if ! docker exec -i "$CONTAINER_ID" /bin/bash < "$1"
then
    # Exit using the variable, to make the build as failure in GitLab CI.
    exit "$BUILD_FAILURE_EXIT_CODE"
fi