#!/bin/bash -eu

setup() {
    cd `dirname $0`
    docker compose up -d
}

setup

docker compose exec go go run src/main.go
