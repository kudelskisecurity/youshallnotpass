#!/usr/bin/env bash
# shellcheck disable=SC1091
# shellcheck disable=SC2206

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

source "${currentDir}/base.sh"

source "${currentDir}/profile.sh"

set -eEo pipefail

# trap any error, and mark it as a system failure.
trap 'exit $SYSTEM_FAILURE_EXIT_CODE' ERR

# Check only the image whitelist on preparation
if [[ -f /usr/local/bin/youshallnotpass ]]; then
    /usr/local/bin/youshallnotpass validate --check-type="image" --ci-platform="gitlab" || exit "$SYSTEM_FAILURE_EXIT_CODE"
elif [[ -f /${currentDir}/youshallnotpass ]]; then
    "${currentDir}/youshallnotpass" validate --check-type="image" --ci-platform="gitlab" || exit "$SYSTEM_FAILURE_EXIT_CODE"
else
    echo "Could Not Find YouShallNotPass Binary"
    exit "$SYSTEM_FAILURE_EXIT_CODE"
fi

wait_for_docker() {
    n=0
    until [ "$n" -ge 5 ]
    do
        docker stats --no-stream >/dev/null && break
        n=$((n+1)) 
        sleep 1
    done
}

is_logged_in() {
  jq -r --arg url "${CUSTOM_ENV_CI_REGISTRY}" '.auths | has($url)' < "$CACHE_DIR/_authfile_$CONTAINER_ID/config.json"
}

start_container() {
    if docker inspect "$CONTAINER_ID" >/dev/null 2>&1; then
        echo 'Found old container, deleting'
        docker kill "$CONTAINER_ID"
        docker rm "$CONTAINER_ID"
    fi

    mkdir -p "$CACHE_DIR/_authfile_$CONTAINER_ID"

    # Use value of ENV variable or {} as empty settings
    echo "${CUSTOM_ENV_DOCKER_AUTH_CONFIG:-{\}}" | jq -r > "$CACHE_DIR/_authfile_$CONTAINER_ID/config.json"
        
    # Try logging into the Gitlab Registry if credentials are provided
    # https://docs.gitlab.com/ee/user/packages/container_registry/index.html#authenticate-by-using-gitlab-cicd
    if [[ "$(is_logged_in)" == "false" ]] && [[ -n "$CUSTOM_ENV_CI_DEPLOY_USER" && -n "$CUSTOM_ENV_CI_DEPLOY_PASSWORD" ]]
    then
        echo "Login to ${CUSTOM_ENV_CI_REGISTRY} with CI_DEPLOY_USER"
        echo "$CUSTOM_ENV_CI_DEPLOY_PASSWORD" | docker --config "$CACHE_DIR/_authfile_$CONTAINER_ID" login  \
            --username "$CUSTOM_ENV_CI_DEPLOY_USER" \
            --password-stdin \
            "$CUSTOM_ENV_CI_REGISTRY" 2>/dev/null
    fi

    if [[ "$(is_logged_in)" == "false" ]] && [[ -n "$CUSTOM_ENV_CI_JOB_USER" && -n "$CUSTOM_ENV_CI_JOB_TOKEN" ]]
    then
        echo "Login to ${CUSTOM_ENV_CI_REGISTRY} with CI_JOB_USER"
        echo "$CUSTOM_ENV_CI_JOB_TOKEN" | docker --config "$CACHE_DIR/_authfile_$CONTAINER_ID" login  \
            --username "$CUSTOM_ENV_CI_JOB_USER" \
            --password-stdin \
            "$CUSTOM_ENV_CI_REGISTRY" 2>/dev/null
    fi

    if [[ "$(is_logged_in)" == "false" ]] && [[ -n "$CUSTOM_ENV_CI_REGISTRY_USER" && -n "$CUSTOM_ENV_CI_REGISTRY_PASSWORD" ]]
    then
        echo "Login to ${CUSTOM_ENV_CI_REGISTRY} with CI_REGISTRY_USER"
        echo "$CUSTOM_ENV_CI_REGISTRY_PASSWORD" | docker --config "$CACHE_DIR/_authfile_$CONTAINER_ID" login  \
            --username "$CUSTOM_ENV_CI_REGISTRY_USER" \
            --password-stdin \
            "$CUSTOM_ENV_CI_REGISTRY" 2>/dev/null
    fi

    # merge default docker config with user-suplied configuration
    if test -f "/root/.docker/config.json"; then
        #cat /root/.docker/config.json "$CACHE_DIR/_authfile_$CONTAINER_ID/config.json" | jq -s '.[0] * .[1]' > "$CACHE_DIR/_authfile_$CONTAINER_ID/config.json" -r
        DOCKER_CONFIG=$(cat "$CACHE_DIR/_authfile_$CONTAINER_ID/config.json" /root/.docker/config.json | jq -sr '.[0] * .[1]')
        echo "$DOCKER_CONFIG" | jq -r > "$CACHE_DIR/_authfile_$CONTAINER_ID/config.json"
    fi

    docker --config="$CACHE_DIR/_authfile_$CONTAINER_ID" pull "$CUSTOM_ENV_CI_JOB_IMAGE"

    rm -rf "$CACHE_DIR/_authfile_$CONTAINER_ID"

    DOCKER_RUN_ARGS_LIST=( $DOCKER_RUN_ARGS )

    docker run \
        --detach \
        --interactive \
        --entrypoint="" \
        --tty \
        --name "$CONTAINER_ID" \
        --volume "$CACHE_DIR:/home/user/cache":Z \
        "${DOCKER_RUN_ARGS_LIST[@]}" \
        "$CUSTOM_ENV_CI_JOB_IMAGE" \
        sleep 999999999
}

install_dependencies() {
    # Copy gitlab-runner binary from the server into the container
    docker cp "$(which gitlab-runner)" "$CONTAINER_ID":/usr/bin/gitlab-runner

    # Install bash in systems with APK (e.g., Alpine)
    docker exec "$CONTAINER_ID" sh -c 'if ! type bash >/dev/null 2>&1 && type apk >/dev/null 2>&1 ; then echo "APK based distro without bash"; apk add bash; fi'

    # Install git in systems with APT (e.g., Debian)
    docker exec "$CONTAINER_ID" /bin/bash -c 'if ! type git >/dev/null 2>&1 && type apt-get >/dev/null 2>&1 ; then echo "APT based distro without git"; apt-get update && apt-get install --no-install-recommends -y ca-certificates git; fi'
    # Install git in systems with DNF (e.g., Fedora)
    docker exec "$CONTAINER_ID" /bin/bash -c 'if ! type git >/dev/null 2>&1 && type dnf >/dev/null 2>&1 ; then echo "DNF based distro without git"; dnf install --setopt=install_weak_deps=False --assumeyes git; fi'
    # Install git in systems with APK (e.g., Alpine)
    docker exec "$CONTAINER_ID" /bin/bash -c 'if ! type git >/dev/null 2>&1 && type apk >/dev/null 2>&1 ; then echo "APK based distro without git"; apk add git; fi'
    # Install git in systems with YUM (e.g., RHEL<=7)
    docker exec "$CONTAINER_ID" /bin/bash -c 'if ! type git >/dev/null 2>&1 && type yum >/dev/null 2>&1 ; then echo "YUM based distro without git"; yum install --assumeyes git; fi'
}

echo "Running in $CONTAINER_ID"

wait_for_docker
start_container
install_dependencies