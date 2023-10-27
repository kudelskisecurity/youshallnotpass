#!/usr/bin/env bash

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" > /dev/null 2>&1 && pwd)"

set -o pipefail

header_template='{
    "type": "JWT",
    "kid": "0001",
    "iss": "github.com"
}'

payload_template='{
    "namespace_id": '"\"${GITHUB_REPOSITORY_OWNER_ID}\""',
    "namespace_path": '"\"${CI_PROJECT_NAMESPACE}\""',
    "project_id": '"\"${GITHUB_REPOSITORY_ID}\""',
    "project_path": '"\"${CI_PROJECT_PATH}\""',
    "user_id": '"\"${GITHUB_ACTOR_ID}\""',
    "user_login": '"\"${GITHUB_TRIGGERING_ACTOR}\""',
    "pipeline_id": '"\"${GITHUB_RUN_ID}\""',
    "job_id": '"\"${GITHUB_RUN_ID}\""',
    "ref": '"\"${GITHUB_REF_NAME}\""',
    "ref_type": '"\"${GITHUB_REF_TYPE}\""',
    "ref_protected": '"\"${GITHUB_REF_PROTECTED}\""',
    "jti": '"\"${ITERATION}\""'
}'

build_header() {
        jq -c \
                --arg iat_str "$(date +%s)" \
                --arg alg "${1:-HS256}" \
        '
        ($iat_str | tonumber) as $iat
        | .alg = $alg
        | .iat = $iat
        | .exp = ($iat + 180)
        ' <<<"$header_template" | tr -d '\n'
}

build_payload() {
        jq -c \
                --arg iat_str "$(date +%s)" \
                --arg project_path "$CI_PROJECT_PATH" \
                --arg namespace_path "$CI_PROJECT_NAMESPACE" \
        '
        ($iat_str | tonumber) as $iat
        | .iat = $iat
        | .exp = ($iat + 180)
        | .project_path = $project_path
        | .namespace_path = $namespace_path
        ' <<<"$payload_template" | tr -d '\n'
}

b64enc() { openssl enc -base64 -A | tr '+/' '-_' | tr -d '='; }
json() { jq -c . | LC_CTYPE=C tr -d '\n'; }
hs_sign() { openssl dgst -binary -sha"${1}" -hmac "$2"; }
rs_sign() { openssl dgst -binary -sha"${1}" -sign <(printf '%s\n' "$2"); }

sign() {
        local algo payload header sig secret=$3
        algo=${1:-RS256}; algo=$(echo "$algo" | tr '[:lower:]' '[:upper:]')
        header=$(build_header "$algo") || return
        payload=${2:-$test_payload}
        signed_content="$(json <<<"$header" | b64enc).$(json <<<"$payload" | b64enc)"
        case $algo in
                HS*) sig=$(printf %s "$signed_content" | hs_sign "${algo#HS}" "$secret" | b64enc) ;;
                RS*) sig=$(printf %s "$signed_content" | rs_sign "${algo#RS}" "$secret" | b64enc) ;;
                *) echo "Unknown algorithm" >&2; return 1 ;;
        esac
        printf '%s.%s\n' "${signed_content}" "${sig}"
}

rsa_secret=$(cat "${currentDir}/certs/private-key.pem")

payload=$(build_payload)
jwt_key=$(sign rs256 "$payload" "$rsa_secret")

echo "$jwt_key"