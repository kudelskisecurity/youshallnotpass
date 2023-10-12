#!/usr/bin/env bash

# Inspired by implementation by Will Haley at:
#   http://willhaley.com/blog/generate-jwt-with-bash/
#   https://stackoverflow.com/questions/46657001/how-do-you-create-an-rs256-jwt-assertion-with-bash-shell-scripting/46672439#46672439

set -o pipefail

# Shared content to use as template
header_template='{
    "typ": "JWT",
    "kid": "0001",
    "iss": "gitlab.example.com"
}'

payload_template='{
  "namespace_id": "1234",
  "namespace_path": "",
  "project_id": "1234",
  "project_path": "",
  "user_id": "123",
  "user_login": "user",
  "user_email": "test.user@example.com",
  "pipeline_id": "1234567",
  "job_id": "123456789",
  "ref": "master",
  "ref_type": "branch",
  "ref_protected": "true",
  "jti": ""
}'

build_header() {
        jq -c \
                --arg iat_str "$(date +%s)" \
                --arg alg "${1:-HS256}" \
        '
        ($iat_str | tonumber) as $iat
        | .alg = $alg
        | .iat = $iat
        | .exp = ($iat + 1)
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
        | .exp = ($iat + 5)
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
        algo=${1:-RS256}; algo=${algo^^}
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

rsa_secret=$(cat /certs/private-key.pem)

payload=$(build_payload)
jwt_key=$(sign rs256 "$payload" "$rsa_secret")

echo "$jwt_key"