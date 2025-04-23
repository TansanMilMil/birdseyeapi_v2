#!/bin/bash -eu

cd `dirname $0`
BASE_DIR="go"
DIST_DIR="$BASE_DIR/dist"
BIN_PATH="$DIST_DIR/birdseyeapi_v2"
GO_FILE="$BASE_DIR/src/main.go"

NO_DOCKER_COMPOSE=false

for ARG in "$@"; do
    case "$ARG" in
        --no-docker-compose)
            NO_DOCKER_COMPOSE=true
            ;;
    esac
done

mkdir -p $DIST_DIR

if $NO_DOCKER_COMPOSE; then
    COMMAND_PREFIX="No docker compose exec go..."
    go build -o $BIN_PATH $GO_FILE
else
    docker compose up -d
    docker compose exec go go build -o $BIN_PATH $GO_FILE
fi

echo "binary built! -> $BIN_PATH"
