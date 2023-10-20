# shellcheck shell=bash disable=SC2034
export YOUSHALLNOTPASS_VAULT_ROOT="cicd"
export VAULT_ROLE="youshallnotpass-github-poc"
export VAULT_LOGIN_PATH="auth/jwt/github.com/login"
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_EXTERNAL_ADDR="http://127.0.0.1:8200"

# for private repos, set those variables
# export GITHUB_USER="your_github_user"
# export GITHUB_TOKEN="your_github_token"