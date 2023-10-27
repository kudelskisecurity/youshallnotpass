#!/usr/bin/env bash
# shellcheck disable=SC1091

TMPDIR=$(pwd)

currentDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

source "${currentDir}/base.sh"

echo "Deleting container $CONTAINER_ID"

docker kill "$CONTAINER_ID"  2>/dev/null || true
docker rm "$CONTAINER_ID"  2>/dev/null || true

# Delete leftover files in /tmp
rm -r "$TMPDIR"

exit 0
