#!/bin/bash -eu

setup() {
    cd `dirname $0`
}

setup
BASE_DIR="go"
DIST_DIR="$BASE_DIR/dist"
BIN_PATH="$DIST_DIR/birdseyeapi_v2"
GO_FILE="$BASE_DIR/src/main.go"

mkdir -p $DIST_DIR

if [ ! -z "${1:-}" ] && [ "$1" == "--no-docker-compose" ]; then
    COMMAND_PREFIX="No docker compose exec go..."
    go build -o $BIN_PATH $GO_FILE
else
    docker compose up -d
    docker compose exec go go build -o $BIN_PATH $GO_FILE
fi

echo "binary built! -> $BIN_PATH"
