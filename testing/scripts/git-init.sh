# shellcheck shell=bash

GIT_REPO="$(pwd)"

today=$(date +%s)

echo "${today}" > "$GIT_REPO/bootstrapped"

cd "$GIT_REPO" || exit

cat >"$GIT_REPO/test.sh" <<EOL
#!/bin/bash

echo "TEST SCRIPT"

EOL

chmod +x "$GIT_REPO/test.sh"

cat >"$GIT_REPO/.gitlab-ci.yml" <<EOL
stages:
  - build

# User MFA Case
user_mfa_job:
  stage: build
  image: alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748
  script:
    - echo "this should require mfa to run"

# User MFA Timeout Case
user_mfa_timeout_job:
  stage: build
  image: alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748
  script:
    - echo "this script should not run because of youshallnotpass timeout"

# Automatic Run Case
automatic_job:
  stage: build
  image: alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748
  script:
    - echo 'this is the automated job and should be run automatically'

# Script Expansion Case
script_job:
  stage: build
  image: alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748
  script:
    - echo "this is the script job"
    - /gitrepo/test.sh

# Failed Execution Case
fail_job:
  stage: build
  image: alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f637
  script:
    - echo 'this job should fail because the image is different'

EOL

git config --global init.defaultBranch master
git config --global user.email "you@example.com"
git config --global user.name "Your Name"
git init
git add .gitlab-ci.yml
git commit -m "commit gitlabci"
git add bootstrapped
git commit -m "commit bootstrapped"
git add test.sh
git commit -m "commit test bash script"
git config remote.origin.url >&- || git remote add origin "$GIT_REPO/.git"
echo "done with git repo initialization time to sleep"
sleep infinity