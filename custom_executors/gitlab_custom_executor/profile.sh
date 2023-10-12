# shellcheck shell=bash disable=SC2034
export YOUSHALLNOTPASS_VAULT_ROOT="cicd"
export VAULT_ROLE="youshallnotpass-demo"
export VAULT_LOGIN_PATH="auth/jwt/gitlab.example.com/login"
export VAULT_ADDR="http://vault:8200"
export VAULT_EXTERNAL_ADDR="http://127.0.0.1:8200"