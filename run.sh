#!/bin/bash -eu

cd `dirname $0`
BASE_DIR="go"
GO_FILE="$BASE_DIR/src/main.go"

NO_DOCKER_COMPOSE=false

for ARG in "$@"; do
    case "$ARG" in
        --no-docker-compose)
            NO_DOCKER_COMPOSE=true
            ;;
    esac
done

if $NO_DOCKER_COMPOSE; then
    COMMAND_PREFIX="No docker compose exec go..."
    go run $GO_FILE
else
    docker compose up -d
    docker compose exec go go run $GO_FILE
fi
