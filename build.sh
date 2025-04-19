#!/bin/bash -eu

setup() {
    cd `dirname $0`
    docker compose up -d
}

setup

docker compose exec go go build -o dist/birdseyeapi_v2 src/main.go
echo 'binary built! -> dist/birdseyeapi_v2'

