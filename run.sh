#!/bin/bash -eu

setup() {
    cd `dirname $0`
}

setup
BASE_DIR="go"
GO_FILE="$BASE_DIR/src/main.go"

if [ ! -z "${1:-}" ] && [ "$1" == "--no-docker-compose" ]; then
    COMMAND_PREFIX="No docker compose exec go..."
    go run $GO_FILE
else
    docker compose up -d
    docker compose exec go go run $GO_FILE
fi
