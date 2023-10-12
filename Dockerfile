FROM golang:1.20@sha256:5865f52f9f277b951610d2ab0b7a14b24cadef7709db26de3320c018fbd4550c AS builder

ARG GOPROXY

WORKDIR /app
COPY . .
RUN go mod download && \
    go mod verity && \
    CGO_ENABLED=0 go build -o main \
    .

FROM gitlab/gitlab-runner:alpine3.18

ARG CI_COMMIT_REF_NAME="development"

ENV RUNNER_BUILDS_DIR="/tmp/builds" \
    RUNNER_CACHE_DIR="/tmp/cache" \
    CUSTOM_CONFIG_EXEC="/var/custom-executor/config.sh" \
    CUSTOM_PREPARE_EXEC="/var/custom-executor/prepare.sh" \
    CUSTOM_RUN_EXEC="/var/custom-executor/run.sh" \
    CUSTOM_CLEANUP_EXEC="/var/custom-executor/cleanup.sh" \
    CUSTOM_CONFIG_EXEC_TIMEOUT=200 \
    CUSTOM_PREPARE_EXEC_TIMEOUT=200 \
    CUSTOM_CLEANUP_EXEC_TIMEOUT=200 \
    CUSTOM_GRACEFUL_KILL_TIMEOUT=200 \
    CUSTOM_FORCE_KILL_TIMEOUT=200

RUN apk add --no-cache docker-cli jq openssl && rm -rf /var/cache/apk/*
RUN mkdir -p /var/custom-executor

COPY --from=builder /app/main /usr/local/bin/youshallnotpass
COPY custom_executors/gitlab_custom_executor/base.sh /var/custom-executor/base.sh
COPY custom_executors/gitlab_custom_executor/cleanup.sh /var/custom-executor/cleanup.sh
COPY custom_executors/gitlab_custom_executor/config.sh /var/custom-executor/config.sh
COPY custom_executors/gitlab_custom_executor/profile.sh /var/custom-executor/profile.sh
COPY custom_executors/gitlab_custom_executor/prepare.sh /var/custom-executor/prepare.sh
COPY custom_executors/gitlab_custom_executor/run.sh /var/custom-executor/run.sh
