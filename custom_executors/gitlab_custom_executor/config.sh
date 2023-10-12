#!/usr/bin/env bash

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

YOUSHALLNOTPASS_PREVALIDATION_TOKEN=$(LC_ALL=C tr -dc '[:alnum:]' < /dev/urandom | head -c50)

if [ "$YOUSHALLNOTPASS_GENERATE_JWT" == true ] && [ -z "${CI_JOB_JWT}" ] && [ -z "${CUSTOM_ENV_CI_JOB_JWT}" ]; then 
    CUSTOM_ENV_CI_JOB_JWT=$(CI_PROJECT_PATH=$CI_PROJECT_PATH CI_NAMESPACE_PATH=$CI_NAMESPACE_PATH "${currentDir}/generate_jwt.sh")
fi

cat << EOS
{
  "driver": {
    "name": "youshallnotpass",
    "version": "${RUNNER_VERSION}"
  },
  "job_env" : {
    "YOUSHALLNOTPASS_PREVALIDATION_TOKEN": "${YOUSHALLNOTPASS_PREVALIDATION_TOKEN}",
    "CUSTOM_ENV_CI_JOB_JWT": "${CUSTOM_ENV_CI_JOB_JWT}"
  }
}
EOS