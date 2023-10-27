#!/bin/sh

if command -v apk && ! command -v jq > /dev/null; then
  echo "Installing jq, openssl..."
  apk update && apk add jq openssl
fi

if [ ! -f "/certs/private-key.pem" ] && [ ! -f "/certs/public-key.pem" ]; then
    echo "Generate Pub/Private keys for JWT tokens"
    openssl genrsa -out /certs/private-key.pem 3072
    openssl rsa -in /certs/private-key.pem -pubout -out /certs/public-key.pem
fi

ROLE="youshallnotpass-demo"
POLICIES="gitlab"

sleep 5s

vault auth enable -path=jwt/gitlab.example.com jwt

vault write auth/jwt/gitlab.example.com/config \
  bond_issuer="gitlab.example.com" \
  default_role="$ROLE" \
  jwt_validation_pubkeys="$(cat /certs/public-key.pem)"


vault write auth/jwt/gitlab.example.com/role/$ROLE -<<EOF
{
  "policies": "gitlab",
  "ttl": "1h",
  "user_claim":"user_login",
  "role_type":"jwt",
  "bound_claims": {
    "namespace_path": "youshallnotpass"
  },
  "claim_mappings": {
    "instance":"iss",
    "namespace_id":"namespace_id",
    "namespace_path":"namespace_path",
    "project_id":"project_id",
    "project_path":"project_path"
  }
}
EOF

if [ -f "policy.hcl" ] ; then
    rm policy.hcl
fi

JWT_ACCESSOR=$(vault auth list -format=json | jq -r '.["jwt/gitlab.example.com/"].accessor')

cat << EOF >> policy.hcl
path "cicd/{{ identity.entity.aliases.$JWT_ACCESSOR.metadata.project_path }}/whitelist" {
  capabilities = ["read", "list"]
}
path "cicd/{{ identity.entity.aliases.$JWT_ACCESSOR.metadata.namespace_path }}/whitelist" {
  capabilities = ["read", "list"]
}
path "cicd/{{ identity.entity.aliases.$JWT_ACCESSOR.metadata.namespace_path }}/youshallnotpass_config" {
  capabilities = ["read", "list"]
}
path "cicd/{{ identity.entity.aliases.$JWT_ACCESSOR.metadata.project_path }}/youshallnotpass_config" {
  capabilities = ["read", "list"]
}
path "cicd/{{ identity.entity.aliases.$JWT_ACCESSOR.metadata.project_path }}/scratch/*" {
  capabilities = ["create", "read", "delete"]
}
EOF

echo "Creating policy"
vault policy write $POLICIES policy.hcl

echo "Enabling KV"

vault secrets enable -path=cicd kv 

printf '{
  "allowed_images":[
      "alpine:3.13@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748"
    ]
  }' | vault kv put cicd/youshallnotpass/whitelist -

printf '{
  "allowed_images":[
      "alpine:3.12.7@sha256:a9c28c813336ece5bb98b36af5b66209ed777a394f4f856c6e62267790883820"
    ],
  "allowed_scripts":[
      "automatic_job@sha256:Ij3eYc5EwfiLD6rPw9qFpN82ydukCduG4bUL9ltQDy4=",
      "script_job@sha256:Kn9ysqTdXVzh52gp2LNiX5RMNRxdoAQytneeLcNsycQ="
    ]
  }' | vault kv put cicd/youshallnotpass/demo/whitelist -

printf '{
  "logger": {
    "name": "console"
  }
}' | vault kv put cicd/youshallnotpass/youshallnotpass_config -

printf '{
  "jobs": [
    {
      "jobName": "user_mfa_job",
      "checks": [
        {
          "name": "mfaRequired",
          "options": {
            "checkType": "script"
          }
        }
      ]
    },
    {
      "jobName": "user_mfa_timeout_job",
      "checks": [
        {
          "name": "imageHash",
          "options": {
            "abortOnFail": false,
            "mfaOnFail": true
          }
        },
        {
          "name": "scriptHash",
          "options": {
            "abortOnFail": false,
            "mfaOnFail": true
          }
        }
      ]
    },
    {
      "jobName": "automatic_job",
      "checks": [
        {
          "name": "imageHash",
          "options": {
            "abortOnFail": true
          }
        },
        {
          "name": "scriptHash",
          "options": {
            "abortOnFail": true
          }
        }
      ]
    },
    {
      "jobName": "script_job",
      "checks": [
        {
          "name": "scriptHash",
          "options": {
            "abortOnFail": true
          }
        }
      ]
    },
    {
      "jobName": "fail_job",
      "checks": [
        {
          "name": "imageHash",
          "options": {
            "abortOnFail": true
          }
        },
        {
          "name": "scriptHash",
          "options": {
            "abortOnFail": true
          }
        }
      ]
    },
    {
      "jobName": "default",
      "checks": [
        {
          "name": "imageHash",
          "options": {
            "abortOnFail": true
          }
        }
      ]
    }
  ]
}
' | vault kv put cicd/youshallnotpass/demo/youshallnotpass_config -